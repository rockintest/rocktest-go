package scenario

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type InlineExecutor struct {
	scenario *Scenario
}

func NewInlineExecutor(s *Scenario) *InlineExecutor {
	ret := new(InlineExecutor)
	ret.scenario = s
	return ret
}

func (x InlineExecutor) Lookup(s string) (string, bool, error) {

	// Do we have expression like
	// ${module(p1,p2).path} or ${module(p1,p2).path::DEFAULT}

	re, err := regexp.Compile(`(?s)\$([^(]+)\(((?:[^,]+)?(?:,[^,]+)*)\)(?:\.([^:]+))?(::){0,1}(.*)`)

	if err != nil {
		log.Errorf(err.Error())
		return "", false, err
	}

	if re.Match([]byte(s)) {
		res := re.FindAllStringSubmatch(s, -1)

		module := res[0][1]
		params := res[0][2]
		path := res[0][3]

		paramArray := strings.Split(params, "]>>,<<[")
		paramArray[0] = strings.TrimPrefix(paramArray[0], "<<[")
		paramArray[len(paramArray)-1] = strings.TrimSuffix(paramArray[len(paramArray)-1], "]>>")

		// Get the Meta information (mainly parameter list)
		meta := x.scenario.Meta(module)

		log.Tracef("Meta informations for module %s: %v", module, meta)

		paramMap := make(map[string]interface{})

		// If no information is available in the Meta
		// use the default:
		//   ext = ext
		//   params = value, param1, param2...
		if meta.Params == nil {
			for i, v := range paramArray {
				if i == 0 {
					paramMap["value"] = v
				} else {
					paramMap[fmt.Sprintf("param%d", i)] = v
				}
			}
		} else {
			for i, v := range meta.Params {
				paramMap[v] = paramArray[i]
			}
		}

		if meta.Ext == "" {
			paramMap["ext"] = path
		} else {
			paramMap[meta.Ext] = path
		}

		err2 := x.scenario.Exec(module, paramMap)

		if err2 != nil {
			log.Errorf("Cannot execute inline step: %s", err2.Error())
			return s, false, err2
		}

		ret, found := x.scenario.Context["??"]

		if found {
			return ret, true, nil
		} else {
			if len(res[0]) > 5 {
				return res[0][5], true, nil
			} else {
				return "", true, nil
			}
		}

	} else {
		return "", false, nil
	}

}
