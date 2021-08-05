package text

import (
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type MapLookup struct {
	Map map[string]string
}

func NewMapLookup(m map[string]string) *MapLookup {
	ret := new(MapLookup)
	ret.Map = m
	return ret
}

func (l MapLookup) getValue(s string) (string, bool) {
	ret, ok := l.Map[s]
	if ok {
		return ret, ok
	} else {
		env := os.Getenv(s)
		if env == "" {
			return s, false
		} else {
			return env, true
		}
	}
}

func (l MapLookup) Lookup(s string) (string, bool, error) {

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

		val, found := l.getValue(name)

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

		str, ok := l.getValue(s)
		return str, ok, nil
	}

}
