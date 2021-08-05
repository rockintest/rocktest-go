package text

import (
	"os"
	"testing"
)

func TestSubst1(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"

	mapLookup := NewMapLookup(themap)

	sub := NewStringSubstitutorByLookuper(mapLookup)

	result, _ := sub.Replace("${rock}")

	expect := "test"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst2(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "rock"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${rock}${test}")

	expect := "testrock"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst3(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "rock"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("$${rock}inter${test}")

	expect := "${rock}interrock"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst4(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "rock"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${rock}inter$${test}")

	expect := "testinter${test}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst5(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "rock"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("$rock}${test}")

	expect := "$rock}rock"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst6(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${rock}${test}")

	expect := "test${test}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst7(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "jazz"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${${rock}}")

	expect := "jazz"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst8(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "jazz"
	themap["testjazz"] = "funk"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${${rock}${test}}")

	expect := "funk"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst9(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "jazz"
	themap["the test and jazz"] = "funk"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${the ${rock} and ${test}}")

	expect := "funk"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst10(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"
	themap["test"] = "jazz"
	themap["jazz"] = "funk"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${${${rock}}}")

	expect := "funk"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst11(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["$rock"] = "test"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${$$rock}")

	expect := "test"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubstEnv1(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${HOME}${rock}")

	expect := os.Getenv("HOME") + "test"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst12(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${rock}\\${ro}ck}")

	expect := "test${ro}ck}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}

func TestSubst14(t *testing.T) {

	var themap map[string]string = make(map[string]string)

	themap["rock"] = "test"

	sub := NewStringSubstitutorByMap(themap)
	result, _ := sub.Replace("${$func(ro\\}ck\\})}")

	expect := "${$func(ro}ck})}"

	if result != expect {
		t.Errorf("Error. Expected %s but was %s", expect, result)
	}

}
