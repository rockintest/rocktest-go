package scenario

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func (module *Module) Regex_matchMeta() Meta {
	return Meta{Ext: "group", Params: []string{"pattern", "string"}}
}

func (module *Module) Regex_match(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	pattern, err := scenario.GetString(paramsEx, "pattern", nil)
	if err != nil {
		return err
	}

	pstring, err := scenario.GetString(paramsEx, "string", nil)
	if err != nil {
		return err
	}

	group, err := scenario.GetNumber(paramsEx, "group", -1)
	if err != nil {
		group = -1
	}

	ml, _ := scenario.GetBool(paramsEx, "multiline", true)

	as, err := scenario.GetString(paramsEx, "as", "match")
	// We have an as parameter, add a "."
	if err == nil && as != "" {
		as = as + "."
	}

	log.Tracef("Regex check. Params=%v", params)

	if ml {
		pattern = "(?s)" + pattern
	}

	re, err := regexp.Compile(pattern)

	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	if re.Match([]byte(pstring)) {

		res := re.FindAllStringSubmatch(pstring, -1)

		log.Tracef("groups: %v", res[0])

		for i, v := range res[0] {
			scenario.PutContext(fmt.Sprintf("%s%d", as, i), v)
		}

		if group == -1 {
			scenario.PutContext(as+"result", true)
			scenario.PutContext("??", true)
		} else {
			if group < len(res[0]) {
				scenario.PutContext(as+"result", res[0][group])
				scenario.PutContext("??", res[0][group])
			} else {
				scenario.PutContext(as+"result", "")
				scenario.DeleteContext("??")
			}
		}

	} else {
		if group == -1 {
			scenario.PutContext(as+"result", false)
			scenario.PutContext("??", false)
		} else {
			scenario.PutContext(as+"result", "")
			scenario.PutContext("??", "")
		}
	}

	return nil
}
