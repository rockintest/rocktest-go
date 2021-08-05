package scenario

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

func (module *Module) jsonGetRoot(str string, path string) (interface{}, error) {

	if strings.HasPrefix(path, "[") {
		return module.jsonGet(str, "$"+path)
	} else {
		return module.jsonGet(str, "$."+path)
	}

}

func (module *Module) jsonGet(str string, path string) (interface{}, error) {
	v := interface{}(nil)

	err := json.Unmarshal([]byte(str), &v)

	if err != nil {
		return nil, err
	}

	ret, err := jsonpath.Get(path, v)
	if err != nil {
		return nil, err
	} else {
		return ret, nil
	}

}

func (module *Module) toJson(src interface{}) (string, error) {

	ret, err := json.Marshal(src)

	if err != nil {
		return "", err
	} else {
		return string(ret), nil
	}

}

func (module *Module) Json_parse(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	path, err := scenario.GetString(paramsEx, "path", nil)
	if err != nil {
		return err
	}

	json, err := scenario.GetString(paramsEx, "json", nil)
	if err != nil {
		return err
	}

	ret1, err := module.jsonGetRoot(json, path)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unknown key") {
			scenario.DeleteContextAs(paramsEx, "parse", "result")
			scenario.DeleteContext("??")
			return nil
		} else {
			return fmt.Errorf(err.Error())
		}
	}

	switch v := ret1.(type) {
	case string:
		scenario.PutContextAs(paramsEx, "parse", "result", v)
		scenario.PutContext("??", v)
	default:
		ret, err := module.toJson(ret1)
		if err != nil {
			return err
		}
		scenario.PutContextAs(paramsEx, "parse", "result", ret)
		scenario.PutContext("??", ret)
	}

	return nil
}

func (module *Module) Json_parseMeta() Meta {
	return Meta{Ext: "path", Params: []string{"json"}}
}
