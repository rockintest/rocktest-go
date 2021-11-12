package scenario

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
	return ret, found, nil

}
