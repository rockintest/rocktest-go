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
		fun := res[0][2]

		log.Debugf("Call function %s", fun)

		steps, ok := scenario.Functions[fun]

		if !ok {
			return fmt.Errorf("function %s does not exist", fun)
		} else {
			scenario.RunSteps(steps)
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
