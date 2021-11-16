package scenario

import (
	"fmt"

	"strings"

	log "github.com/sirupsen/logrus"
)

type Step struct {
	Origin map[string]interface{}
	Type   string
	Desc   string
	Value  string
	Params map[string]interface{}
	Steps  []interface{}

	M        Module
	scenario *Scenario
}

// Construct a new Step from the YAML node
func NewStep(n map[string]interface{}, s *Scenario) *Step {

	ret := new(Step)

	ret.Origin = n
	ret.scenario = s

	var val interface{}
	var as string

	for k, v := range n {

		log.Tracef("%s - %s", k, v)

		switch k {
		case "desc":
			ret.Desc = fmt.Sprintf("%v", v)
		case "params":
			ret.Params = nodeToMap(v)
		case "steps":
			ret.Steps = nodeToList(v)
		case "as":
			as = fmt.Sprint(v)
		default:
			ret.Type = strings.Title(strings.ReplaceAll(k, ".", "_"))
			val = v
			if v != nil {
				ret.Value = fmt.Sprintf("%v", v)
			} else {
				ret.Value = ""
			}
		}
	}

	if ret.Params == nil {
		ret.Params = make(map[string]interface{})
	}

	if val != nil {
		ret.Params["value"] = val
	}

	if as != "" {
		ret.Params["as"] = as
	}

	return ret

}

func (s Step) ToString() string {
	ret := fmt.Sprintf("Type: %s, Value: %s, Desc: %s, Params: %v", s.Type, s.Value, s.Desc, s.Params)
	return ret
}

func (s *Step) Exec() error {

	log.Tracef("STEP: %s", s.ToString())
	if s.Desc != "" {
		fmt.Printf(">> %s\n", s.Desc)
	}

	return s.scenario.Exec(s.Type, s.Params)

}
