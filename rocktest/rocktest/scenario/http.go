package scenario

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Handler struct {
	scenario *Scenario
	params   map[string]interface{}
	server   *http.Server
}

func (handler *Handler) findCondition(conditions []interface{}, req *http.Request) (map[string]interface{}, error) {

	for _, cond := range conditions {

		// The type is already checked. We can cast safely
		condMap := cond.(map[string]interface{})

		regexUri, _ := handler.scenario.GetString(condMap, "uri", "")
		regexMeth, _ := handler.scenario.GetString(condMap, "method", "")

		regexMeth = strings.ToUpper(regexMeth)

		// Expand the variables in the URI and in the method
		regexUri = handler.scenario.ExpandString(regexUri)
		regexMeth = handler.scenario.ExpandString(regexMeth)

		reUri, err := regexp.Compile(regexUri)
		if err != nil {
			return nil, err
		}

		reMeth, err := regexp.Compile(regexMeth)
		if err != nil {
			return nil, err
		}

		if (regexUri == "" || reUri.Match([]byte(req.RequestURI))) &&
			(regexMeth == "" || reMeth.Match([]byte(req.Method))) {

			return condMap, nil

		}
	}

	return nil, nil

}

func (handler *Handler) handleRequest(w http.ResponseWriter, req *http.Request) {

	log.Infof("Request %s - URI = %s", req.Method, req.RequestURI)

	if req.RequestURI == "/shutdown" {

		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	req.Body.Close()
	body := buf.String()

	log.Debugf("Body:\n%v", body)

	handler.scenario.PutContext("body", body)
	handler.scenario.PutContext("uri", req.RequestURI)
	handler.scenario.PutContext("method", req.Method)

	conditions, err := handler.scenario.GetList(handler.params, "when", nil)

	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	cond, err := handler.findCondition(conditions, req)

	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if cond != nil {

		// Remove all previous groups
		handler.scenario.DeleteContextRegex("[0-9]+")
		handler.scenario.DeleteContextRegex("urlpath.[0-9]+")
		handler.scenario.DeleteContextRegex("urlvar.[0-9]+")

		// Populate variables urlpath.XX with path elements
		for i, v := range strings.Split(req.URL.Path, "/") {
			handler.scenario.PutContext(fmt.Sprintf("urlpath.%d", i), v)
		}

		// Populate variables urlvar.XX with the variables
		for k, v := range req.URL.Query() {
			handler.scenario.PutContext(fmt.Sprintf("urlvar.%s", k), v[0])
		}

		// Add global headers, if any
		globalHeaders, _ := handler.scenario.GetMap(handler.params, "headers", nil)
		globalHeaders = handler.scenario.ExpandMap(globalHeaders)
		for k, v := range globalHeaders {
			w.Header().Set(k, fmt.Sprint(v))
		}

		// Do the group extraction using the regex
		regexUri, _ := handler.scenario.GetString(cond, "uri", "")
		if regexUri != "" {

			reUri, err := regexp.Compile(handler.scenario.ExpandString(regexUri))
			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			groups := reUri.FindAllStringSubmatch(req.RequestURI, -1)

			if groups == nil {
				http.Error(w, "Internal server error: Should match", http.StatusInternalServerError)
				return
			}

			for i, v := range groups[0] {
				handler.scenario.PutContext(fmt.Sprint(i), v)
			}

		}

		condEx := handler.scenario.ExpandMap(cond)

		// If steps are attached to the condition, execute them
		steps, _ := handler.scenario.GetList(condEx, "steps", nil)
		if steps != nil {
			err := handler.scenario.RunSteps(steps)
			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			// Expend condition again, to replace the variables declared in the steps
			condEx = handler.scenario.ExpandMap(condEx)
		}

		respString, err := handler.scenario.GetString(condEx, "response", nil)

		// reponse field is a string. It is the body, and we return a 200 code
		if err == nil {
			fmt.Fprint(w, respString)
		} else {

			respMap, err := handler.scenario.GetMap(condEx, "response", nil)

			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			// Add header specific to each condition, if any
			localHeaders, _ := handler.scenario.GetMap(respMap, "headers", nil)
			for k, v := range localHeaders {
				w.Header().Set(k, fmt.Sprint(v))
			}

			code, _ := handler.scenario.GetNumber(respMap, "code", 200)
			body, _ := handler.scenario.GetString(respMap, "body", "")

			w.WriteHeader(code)
			fmt.Fprint(w, body)

		}

	} else {
		http.Error(w, fmt.Sprintf("No match for URI %s with method %s", req.RequestURI, req.Method), http.StatusNotFound)
	}
}

func serve(h *Handler, port int) error {

	err := h.server.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Errorf("Error runing mock: %v", err)
		h.scenario.ErrorChan <- err
		return err
	} else {
		log.Info("Mock stopped")
		return nil
	}
}

func (module *Module) checkConditions(conditions []interface{}) error {
	for _, cond := range conditions {

		switch cond.(type) {
		case map[string]interface{}:
		default:
			return fmt.Errorf("bad type for condition %v. Must be a map, not %T", cond, cond)
		}

	}
	return nil
}

func (module *Module) Http_mock(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)
	port, err := scenario.GetNumber(paramsEx, "port", nil)

	if err != nil {
		return err
	}

	conditions, err := scenario.GetList(params, "when", nil)

	if err != nil {
		return err
	}

	err = module.checkConditions(conditions)
	if err != nil {
		return err
	}

	h := &Handler{scenario: scenario, params: params}

	log.Infof("Start Mock on port %d", port)

	block, _ := scenario.GetBool(params, "block", false)

	m := http.NewServeMux()
	s := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      m,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	h.server = &s
	m.HandleFunc("/", h.handleRequest)

	if block {
		serve(h, port)
	} else {
		scenario.PutStore("mockserver", &s)
		scenario.PutCleanup("mockserver", shutdownMock)
		go serve(h, port)

		// Wait a little, then check that we do not have an error
		time.Sleep(time.Duration(200) * time.Millisecond)

		select {
		case err := <-scenario.ErrorChan:
			return err
		default:
		}

	}
	return nil

}

// Module to shutdown current running mock
func (module *Module) Http_shutdownmock(params map[string]interface{}, scenario *Scenario) error {
	return shutdownMock(scenario)
}

// Check conditions
func (module *Module) check(expect map[string]interface{}, code int, body string, scenario *Scenario) error {

	expectedCode, err := scenario.GetNumber(expect, "code", nil)

	if err == nil {

		log.Debugf("Check return code. Expected = %d, actual = %d", expectedCode, code)

		if expectedCode != code {
			return fmt.Errorf("bad HTTP result code. Expected %d but was %d", expectedCode, code)
		}

	}

	jsonMap, err := scenario.GetMap(expect, "body.json", nil)

	if err == nil {

		for k, v := range jsonMap {
			val, err := module.jsonGet(body, "$."+k)

			var actual string

			if err != nil {
				return err
			}

			switch s := val.(type) {
			case string:
				actual = s
			default:
				actual, err = module.toJson(val)
				if err != nil {
					return err
				}
			}

			if actual != v {
				return fmt.Errorf("no matching path=%s, expected %s but was %s", k, v, actual)
			}

		}

	}

	return nil
}

// Module do to HTTP Get requests
func (module *Module) Http_get(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)
	url, err := scenario.GetString(paramsEx, "url", nil)

	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	scenario.PutContextAs(paramsEx, "get", "body", string(body))
	scenario.PutContextAs(paramsEx, "get", "code", resp.StatusCode)

	expect, err := scenario.GetMap(paramsEx, "expect", nil)

	if err == nil {
		return module.check(expect, resp.StatusCode, string(body), scenario)
	}

	return nil
}

func shutdownMock(scenario *Scenario) error {
	log.Info("Cleanup mocks")
	srv := scenario.GetStore("mockserver")

	if srv != nil {
		scenario.RemoveStore("mockserver")
		srv2 := srv.(*http.Server)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		srv2.SetKeepAlivesEnabled(false)
		ret := srv2.Shutdown(ctx)
		return ret

	}
	return nil
}
