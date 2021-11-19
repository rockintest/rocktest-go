package scenario

func (module *Module) Loop(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	from, _ := scenario.GetNumber(paramsEx, "from", 0)
	inc, _ := scenario.GetNumber(paramsEx, "inc", 1)
	counter, _ := scenario.GetString(paramsEx, "counter", "i")
	to, err := scenario.GetNumber(paramsEx, "to", nil)
	if err != nil {
		return err
	}

	steps, _ := scenario.GetList(params, "steps", nil)
	if err != nil {
		return err
	}

	for i := from; i < to; i += inc {
		scenario.PutContext(counter, i)
		scenario.RunSteps(steps)
	}

	return nil
}
