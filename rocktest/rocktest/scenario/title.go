package scenario

import (
	"fmt"
)

func (module *Module) Title(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)

	val, _ := scenario.GetString(paramsEx, "value", "")
	fmt.Print("========================\n")
	fmt.Printf("==     %s\n", val)
	fmt.Print("========================\n")

	return nil
}
