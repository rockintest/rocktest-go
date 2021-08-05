package scenario

import "strings"

func (module *Module) Toupper(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	val, _ := scenario.GetString(paramsEx, "value", "")
	as, _ := scenario.GetString(paramsEx, "as", "toupper.result")

	ret := strings.ToUpper(val)

	scenario.PutContext(as, ret)
	scenario.PutContext("??", ret)

	return nil
}

func (module *Module) Tolower(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	val, _ := scenario.GetString(paramsEx, "value", "")
	as, _ := scenario.GetString(paramsEx, "as", "tolower.result")

	ret := strings.ToLower(val)

	scenario.PutContext(as, ret)
	scenario.PutContext("??", ret)

	return nil
}
