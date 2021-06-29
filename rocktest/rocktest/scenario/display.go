package scenario

import (
	"fmt"
)

func (module *Module) Display(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)

	val, _ := scenario.GetString(paramsEx, "value", "")
	fmt.Printf(">> %s\n", val)

	return nil
}
