package scenario

import (
	"fmt"
)

func (module *Module) Display(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	val, _ := scenario.GetString(paramsEx, "value", "")
	fmt.Printf(">> %s\n", val)

	return nil
}
