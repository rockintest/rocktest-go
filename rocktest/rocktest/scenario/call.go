package scenario

import (
	"fmt"
	"strings"
)

func (module *Module) Call(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)

	val, err := scenario.GetString(paramsEx, "value", nil)

	if err != nil {
		return err
	}

	calledScenario := NewScenario()

	calledScenario.Caller = scenario
	calledScenario.Root = scenario.Root

	if !strings.HasSuffix(val, ".yaml") {
		val = val + ".yaml"
	}

	contextStr, err := scenario.GetString(paramsEx, "context", nil)
	if err == nil {
		// We have something like context: XXX in the params
		if contextStr == "all" {
			calledScenario.CopyVariables(scenario)
		} else {
			calledScenario.CopyVariable(contextStr, scenario)
		}
	} else {
		// Check if we have a list in context param
		contextList, err := scenario.GetList(paramsEx, "context", nil)
		if err == nil {
			for _, v := range contextList {
				calledScenario.CopyVariable(fmt.Sprintf("%v", v), scenario)
			}
		}
	}

	calledScenario.AddVariables(paramsEx)

	err = calledScenario.RunFromRoot(val)

	if err != nil {
		return err
	}

	return nil
}
