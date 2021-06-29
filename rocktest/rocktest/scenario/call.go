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

	context, err := scenario.GetList(paramsEx, "context", nil)
	if err == nil {
		if context[0] == "all" {
			calledScenario.CopyVariables(scenario)
		} else {
			for _, v := range context {
				calledScenario.CopyVariable(fmt.Sprint(v), scenario)
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
