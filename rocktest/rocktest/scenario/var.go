package scenario

import (
	"errors"
	"regexp"
)

func (module *Module) Var(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)

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
		scenario.PutContext(res[0][1], res[0][2])
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
		return errors.New("Variable delaration must be VAR=VALUE, not " + val + ". Or you can use a params map with name & value fields")
	}

	scenario.PutContext(name, val)

	return nil
}
