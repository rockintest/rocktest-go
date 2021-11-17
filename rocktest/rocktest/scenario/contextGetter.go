package scenario

import (
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type ContextGetter struct {
	scenario *Scenario
}

func NewContextGetter(s *Scenario) *ContextGetter {
	ret := new(ContextGetter)
	ret.scenario = s
	return ret
}

func (x ContextGetter) Lookup(s string) (string, bool, error) {

	// Do we have expression like
	// ${variable?value if set::value if not set} or
	// ${variable::value if not set}

	re, err := regexp.Compile(`([^?]+)(?:\?(.*))?::(.*)`)

	if err != nil {
		log.Errorf(err.Error())
		return s, false, nil
	}

	if re.Match([]byte(s)) {
		res := re.FindAllStringSubmatch(s, -1)

		name := res[0][1]
		valset := res[0][2]
		valunset := res[0][3]

		log.Tracef("%v", res)
		log.Tracef("Name: \"%s\", Val is set: \"%s\", val if not set: \"%s\"", name, valset, valunset)

		val, found := x.scenario.GetContext(name)

		if found {
			if valset != "" {
				return valset, true, nil
			} else {
				return val, true, nil
			}
		} else {
			if valunset != "" {
				return valunset, true, nil
			} else {
				return s, false, nil
			}
		}

	} else {

		ret, found := x.scenario.GetContext(s)

		if !found {
			env := os.Getenv(s)
			if env == "" {
				return s, false, nil
			} else {
				return env, true, nil
			}

		} else {
			return ret, found, nil
		}
	}

}
