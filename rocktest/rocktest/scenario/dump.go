package scenario

import (
	"fmt"
)

func (module *Module) Dump(params map[string]interface{}, scenario *Scenario) error {

	mod, _ := scenario.GetContext("module")

	fmt.Printf("Variables for context %s\n", mod)

	for k, v := range scenario.getCurrentContext() {
		fmt.Printf("  %s = %v\n", k, v)
	}

	return nil
}
