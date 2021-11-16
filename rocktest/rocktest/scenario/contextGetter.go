package scenario

import "os"

type ContextGetter struct {
	scenario *Scenario
}

func NewContextGetter(s *Scenario) *ContextGetter {
	ret := new(ContextGetter)
	ret.scenario = s
	return ret
}

func (x ContextGetter) Lookup(s string) (string, bool, error) {

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
