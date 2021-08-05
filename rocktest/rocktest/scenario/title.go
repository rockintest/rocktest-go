package scenario

import (
	"fmt"
)

func (module *Module) Title(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	val, _ := scenario.GetString(paramsEx, "value", "")
	fmt.Print("========================\n")
	fmt.Printf("==     %s\n", val)
	fmt.Print("========================\n")

	return nil
}
