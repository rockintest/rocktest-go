package scenario

import (
	"errors"
	"regexp"
	"strings"
)

func putCaller(name string, val string, scenario *Scenario) {

	if scenario.InFunction {

		if strings.HasPrefix(name, ".") {
			scenario.PutContextCallerFunctio(name[1:], val)
		} else {
			scenario.PutContextCallerFunctio(scenario.GetModule()+"."+name, val)
		}

	} else {

		if strings.HasPrefix(name, ".") {
			scenario.Caller.PutContext(name[1:], val)
		} else {
			scenario.Caller.PutContext(scenario.GetModule()+"."+name, val)
		}

	}
}

func (module *Module) Return(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	val, err := scenario.GetString(paramsEx, "value", nil)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(` *([^ ]*) *= *(.*) *`)

	if err != nil {
		return err
	}

	// The variable is declared as VAR = VALUE
	// We have something like :
	// - var: VARNAME = VALUE
	if re.Match([]byte(val)) {
		res := re.FindAllStringSubmatch(val, -1)
		putCaller(res[0][1], res[0][2], scenario)
		return nil
	}

	// We have something like
	// - var:
	//   params:
	//     name: VARNAME
	//     value: VALUE
	// Which is handy if the value is on multiple lines
	// In this case, we need value AND name fields

	name, err := scenario.GetString(paramsEx, "name", nil)
	if err != nil {
		return errors.New("Variable return must be VAR=VALUE, not " + val + ". Or you can use a params map with name & value fields")
	}

	putCaller(name, val, scenario)

	return nil
}
