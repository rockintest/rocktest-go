package scenario

// Skips execution of next steps
func (module *Module) Skip(params map[string]interface{}, scenario *Scenario) error {

	scenario.Skip = true
	return nil

}

// Resumes execution
func (module *Module) Resume(params map[string]interface{}, scenario *Scenario) error {

	scenario.Skip = false
	return nil

}
