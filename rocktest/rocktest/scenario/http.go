package scenario

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Handler struct {
	scenario *Scenario
	params   map[string]interface{}
	server   http.Server
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

func serve(h *Handler, port int) {

	//m := http.NewServeMux()
	//s := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: m}
	//h.server = s
	//m.HandleFunc("/", h.handleRequest)

	h.server.ListenAndServe()
	log.Info("Mock stopped")

	//http.HandleFunc("/", h.handleRequest)
	//http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  1 * time.Second,
	}

	h.server = s
	m.HandleFunc("/", h.handleRequest)

	if block {
		serve(h, port)
	} else {
		scenario.PutStore("mockserver", s)
		go serve(h, port)

		//ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		//defer cancel()

		//s.Shutdown(ctx)
	}
	return nil

}

func (module *Module) Http_shutdownmock(params map[string]interface{}, scenario *Scenario) error {
	srv := scenario.GetStore("mockserver")

	if srv != nil {
		scenario.RemoveStore("mockserver")
		srv2 := srv.(http.Server)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		srv2.SetKeepAlivesEnabled(false)
		ret := srv2.Shutdown(ctx)
		return ret
	}

	return nil
}
