package scenario

import (
	"fmt"
)

func (module *Module) Dump(params map[string]interface{}, scenario *Scenario) error {

	fmt.Printf("Variables for context %s\n", scenario.Context["module"])

	for k, v := range scenario.Context {
		fmt.Printf("  %s = %v\n", k, v)
	}

	return nil
}
