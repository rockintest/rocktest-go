package scenario

import (
	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
)

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
