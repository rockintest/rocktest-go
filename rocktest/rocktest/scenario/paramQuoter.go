package scenario

import (
	"bytes"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ParamQuoter struct {
	Scen *Scenario
}

func (l ParamQuoter) Lookup(s string) (string, bool, error) {

	// Do we have expression like
	// ${module(p1,p2).path}

	re, err := regexp.Compile(`(?s)\$([^(]+)\(((?:[^,]+)?(?:,[^,]+)*)\)(?:\.(.+))?`)

	if err != nil {
		log.Errorf(err.Error())
		return "", false, nil
	}

	if re.Match([]byte(s)) {
		res := re.FindAllStringSubmatch(s, -1)

		module := res[0][1]
		params := res[0][2]
		path := res[0][3]

		log.Tracef("%s - %s - %s", module, params, path)

		if strings.HasPrefix(params, "<<[") || params == "" {
			return "", false, nil
		}

		var ret bytes.Buffer
		ret.WriteString("${$")
		ret.WriteString(module)
		ret.WriteString("(")

		paramArray := strings.Split(params, ",")

		for i, v := range paramArray {
			ret.WriteString("<<[")

			vSubst, _ := l.Scen.Subst.Replace(v)
			vSubst = strings.ReplaceAll(vSubst, "{", "\\{")
			vSubst = strings.ReplaceAll(vSubst, "}", "\\}")

			ret.WriteString(vSubst)
			ret.WriteString("]>>")
			if i != len(paramArray)-1 {
				ret.WriteString(",")
			}
		}

		ret.WriteString(")")

		if path != "" {
			ret.WriteString(".")
			ret.WriteString(path)
		}

		ret.WriteString("}")

		return ret.String(), true, nil
	}

	return "", false, nil
}
