package scenario

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (module *Module) Call(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	val, err := scenario.GetString(paramsEx, "value", nil)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(` *(.*)->(.*) *`)

	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	if re.Match([]byte(val)) {

		// Call a function

		res := re.FindAllStringSubmatch(val, -1)
		mod := res[0][1]
		fun := res[0][2]

		if mod == "" {
			log.Debugf("Call local function %s", fun)

			steps, ok := scenario.Functions[fun]

			if !ok {
				return fmt.Errorf("function %s does not exist", fun)
			} else {

				inFunctionBefore := scenario.InFunction
				scenario.InFunction = true
				scenario.pushContext()

				// The module name is actually the function name for this call
				scenario.PutContext("module", fun)
				scenario.AddVariables(paramsEx)
				err := scenario.RunSteps(steps)
				if err != nil {
					return err
				}
				scenario.popContext()
				scenario.InFunction = inFunctionBefore

			}
		} else {
			log.Debugf("Call function %s in module %d", fun, mod)

			calledScenario := NewScenario()

			calledScenario.Caller = scenario
			calledScenario.Root = scenario.Root

			if !strings.HasSuffix(mod, ".yaml") {
				mod = mod + ".yaml"
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

			err = calledScenario.RunFunction(mod, fun)
			if err != nil {
				return err
			}
		}

	} else {

		// Call a scenario

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

	}

	return nil
}
