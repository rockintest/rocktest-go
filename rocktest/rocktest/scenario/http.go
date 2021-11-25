package scenario

import (
	"bytes"
	"context"
	"crypto/tls"
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

		regexUri, err := handler.scenario.GetString(condMap, "uri", "")
		if err != nil {
			return nil, err
		}
		regexMeth, err := handler.scenario.GetString(condMap, "method", "")
		if err != nil {
			return nil, err
		}

		regexMeth = strings.ToUpper(regexMeth)

		// Expand the variables in the URI and in the method
		regexUri, err = handler.scenario.ExpandString(regexUri)
		if err != nil {
			return nil, err
		}

		regexMeth, err = handler.scenario.ExpandString(regexMeth)
		if err != nil {
			return nil, err
		}

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
		handler.scenario.DeleteContextRegex("headers.[0-9]+")

		// Populate variables urlpath.XX with path elements
		for i, v := range strings.Split(req.URL.Path, "/") {
			handler.scenario.PutContext(fmt.Sprintf("urlpath.%d", i), v)
		}

		// Populate variables urlvar.XX with the variables
		for k, v := range req.URL.Query() {
			handler.scenario.PutContext(fmt.Sprintf("urlvar.%s", k), v[0])
		}

		// Populate variables headers.XX with the headers
		for k, v := range req.Header {
			handler.scenario.PutContext(fmt.Sprintf("headers.%s", strings.ToLower(k)), v[0])
		}

		// Add global headers, if any
		globalHeaders, _ := handler.scenario.GetMap(handler.params, "headers", nil)
		globalHeaders, err = handler.scenario.ExpandMap(globalHeaders)
		if err != nil {
			http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		for k, v := range globalHeaders {
			w.Header().Set(k, fmt.Sprint(v))
		}

		// Do the group extraction using the regex
		regexUri, _ := handler.scenario.GetString(cond, "uri", "")
		if regexUri != "" {

			val, err := handler.scenario.ExpandString(regexUri)
			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
				return
			}

			reUri, err := regexp.Compile(val)
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

		condEx, err := handler.scenario.ExpandMap(cond)
		if err != nil {
			http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		// If steps are attached to the condition, execute them
		steps, _ := handler.scenario.GetList(condEx, "steps", nil)
		if steps != nil {
			err := handler.scenario.RunSteps(steps)
			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			// Expend condition again, to replace the variables declared in the steps
			condEx, err = handler.scenario.ExpandMap(condEx)
			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
				return
			}

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

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

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
func (module *Module) check(expect map[string]interface{}, code int, body string, h http.Header, scenario *Scenario) error {

	expectedCode, err := scenario.GetString(expect, "code", nil)

	if err == nil {

		log.Debugf("Check return code. Expected = %s, actual = %d", expectedCode, code)

		ok, err := regexp.MatchString(expectedCode, fmt.Sprint(code))

		if err != nil {
			return err
		}

		if !ok {
			return fmt.Errorf("bad HTTP result code. Expected %s but was %d", expectedCode, code)
		}

	}

	jsonMap, err := scenario.GetMap(expect, "body.json", nil)

	if err == nil {

		for k, v := range jsonMap {
			val, err := module.jsonGetRoot(body, k)

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

			ok, err := regexp.MatchString(fmt.Sprint(v), actual)

			if err != nil {
				return err
			}

			if !ok {
				log.Debugf("Check if regex '%s' match '%s' => NO", fmt.Sprint(v), actual)
				return fmt.Errorf("no matching path=%s, expected %s but was %s", k, v, actual)
			} else {
				log.Debugf("Check if regex '%s' match '%s' => YES", fmt.Sprint(v), actual)
			}

		}

	}

	exprs, err := scenario.GetList(expect, "body.match", nil)

	if err == nil {
		for _, expr := range exprs {
			ok, err := regexp.MatchString(fmt.Sprintf("(?s)%v", expr), body)

			if err != nil {
				return err
			}

			if !ok {
				log.Debugf("Check if regex '%v' match the body => NO", expr)
				return fmt.Errorf("body not matching '%v'", expr)
			} else {
				log.Debugf("Check if regex '%v' match the body => YES", expr)
			}

		}
	}

	headers, err := scenario.GetMap(expect, "headers", nil)

	if err == nil {

		for k, v := range headers {

			ok, err := regexp.MatchString(fmt.Sprint(v), h.Get(k))

			if err != nil {
				return err
			}

			if !ok {
				log.Debugf("Check if header %s match regex '%s' => NO", k, fmt.Sprint(v))
				return fmt.Errorf("header %s does not match regex. Expected %v but was %s", k, v, h.Get(k))
			} else {
				log.Debugf("Check if header %s match regex '%s' => YES", k, fmt.Sprint(v))
			}

		}

	}

	return nil
}

// Do a http request
func (module *Module) httpReq(params map[string]interface{}, meth string, scenario *Scenario) error {

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	url, err := scenario.GetString(paramsEx, "url", nil)

	method := strings.ToLower(meth)

	if err != nil {
		return err
	}

	bodyRequest, err := scenario.GetString(paramsEx, "body", nil)

	var bodyReader io.Reader = nil
	if err == nil {
		bodyReader = strings.NewReader(bodyRequest)
	}

	client := &http.Client{}
	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)

	if err != nil {
		return err
	}

	headers, err := scenario.GetMap(paramsEx, "headers", nil)

	if err == nil {
		for k, v := range headers {
			req.Header.Add(k, fmt.Sprint(v))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyResponse, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	scenario.PutContextAs(paramsEx, method, "body", string(bodyResponse))
	scenario.PutContextAs(paramsEx, method, "code", resp.StatusCode)

	for k, v := range resp.Header {
		scenario.PutContextAs(paramsEx, method, "headers."+strings.ToLower(k), v[0])
	}

	expect, err := scenario.GetMap(paramsEx, "expect", nil)

	if err == nil {
		return module.check(expect, resp.StatusCode, string(bodyResponse), resp.Header, scenario)
	}

	return nil

}

// Module do to HTTP Get requests
func (module *Module) Http_get(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "get", scenario)

}

// Module do to HTTP Post requests
func (module *Module) Http_post(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "post", scenario)

}

// Module do to HTTP Put requests
func (module *Module) Http_put(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "put", scenario)

}

// Module do to HTTP Head requests
func (module *Module) Http_head(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "head", scenario)

}

// Module do to HTTP Delete requests
func (module *Module) Http_delete(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "delete", scenario)

}

// Module do to HTTP Connect requests
func (module *Module) Http_connect(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "connect", scenario)

}

// Module do to HTTP Options requests
func (module *Module) Http_options(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "options", scenario)

}

// Module do to HTTP Trace requests
func (module *Module) Http_trace(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "trace", scenario)

}

// Module do to HTTP Patch requests
func (module *Module) Http_patch(params map[string]interface{}, scenario *Scenario) error {

	return module.httpReq(params, "patch", scenario)

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
