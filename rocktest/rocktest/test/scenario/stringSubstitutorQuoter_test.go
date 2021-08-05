package scenarioTest

import (
	"testing"

	"io.rocktest/rocktest/scenario"
	"io.rocktest/rocktest/text"
)

func TestSubstQuoter1(t *testing.T) {

	var q scenario.ParamQuoter
	sc := scenario.NewScenario()
	q.Scen = sc

	sub := text.NewStringSubstitutorByLookuper(q)

	result, _ := sub.Replace("${rock}")

	expect := "${rock}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubstQuoter2(t *testing.T) {

	var q scenario.ParamQuoter
	sc := scenario.NewScenario()
	q.Scen = sc

	sub := text.NewStringSubstitutorByLookuper(q)

	result, _ := sub.Replace("${$rock()}")

	expect := "${$rock()}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubstQuoter3(t *testing.T) {

	var q scenario.ParamQuoter
	sc := scenario.NewScenario()
	q.Scen = sc

	sub := text.NewStringSubstitutorByLookuper(q)

	result, _ := sub.Replace("${$rock(param)}")

	expect := "${$rock(<<[param]>>)}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubstQuoter4(t *testing.T) {

	var q scenario.ParamQuoter
	sc := scenario.NewScenario()
	q.Scen = sc

	sub := text.NewStringSubstitutorByLookuper(q)

	result, _ := sub.Replace("${$rock(param,param2)}")

	expect := "${$rock(<<[param]>>,<<[param2]>>)}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubstQuoter5(t *testing.T) {

	var q scenario.ParamQuoter
	sc := scenario.NewScenario()
	q.Scen = sc

	sub := text.NewStringSubstitutorByLookuper(q)

	result, _ := sub.Replace("${$rock(param,param2).path}")

	expect := "${$rock(<<[param]>>,<<[param2]>>).path}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubstQuoter6(t *testing.T) {

	var q scenario.ParamQuoter
	sc := scenario.NewScenario()
	q.Scen = sc

	sub := text.NewStringSubstitutorByLookuper(q)

	result, _ := sub.Replace("${$rock(<<[param]>>)}")

	expect := "${$rock(<<[param]>>)}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubstQuoter7(t *testing.T) {

	var q scenario.ParamQuoter
	sc := scenario.NewScenario()
	q.Scen = sc

	sub := text.NewStringSubstitutorByLookuper(q)

	result, _ := sub.Replace("${${rock}}")

	expect := "${${rock}}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}
