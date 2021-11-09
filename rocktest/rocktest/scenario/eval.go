package scenario

import (
	"fmt"
)

func (module *Module) Eval_evalMeta() Meta {
	return Meta{Ext: "path", Params: []string{"expr"}}
}

func (module *Module) Eval(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	expr, _ := scenario.GetString(paramsEx, "expr", "")
	fmt.Printf(">> %s\n", expr)

	return nil
}
