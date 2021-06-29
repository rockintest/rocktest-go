package scenarioTest

import (
	"testing"

	"io.rocktest/rocktest/scenario"
)

func TestGetString1(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = "test"

	s := scenario.NewScenario()

	str, err := s.GetString(m, "rock", nil)

	if str != "test" {
		t.Errorf("Bad result. Expected %s but was %s", "rock", str)
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}

func TestGetString2(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = "test"

	s := scenario.NewScenario()

	str, err := s.GetString(m, "test", "rock")

	if str != "rock" {
		t.Errorf("Bad result. Expected %s but was %s", "test", str)
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}

func TestGetString3(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = "test"

	s := scenario.NewScenario()

	_, err := s.GetString(m, "test", nil)

	if err == nil {
		t.Errorf("Should return an error")
	}

	t.Logf("Error returned: %s", err.Error())
}

func TestGetString4(t *testing.T) {

	s := scenario.NewScenario()

	str, err := s.GetString(nil, "test", "rock")

	if str != "rock" {
		t.Errorf("Bad result. Expected %s but was %s", "test", str)
	}

	if err != nil {
		t.Errorf("Should not return an error: %s", err.Error())
	}

}

func TestGetString5(t *testing.T) {

	s := scenario.NewScenario()

	_, err := s.GetString(nil, "test", nil)

	if err == nil {
		t.Errorf("Should return an error")
	}

	t.Logf("Error returned: %s", err.Error())
}

func TestGetString6(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = 12

	s := scenario.NewScenario()

	str, err := s.GetString(m, "rock", nil)

	if str != "12" {
		t.Errorf("Bad result. Expected %s but was %s", "12", str)
	}

	if err != nil {
		t.Errorf("Should not return an error: %s", err.Error())
	}

}

//------

func TestGetList1(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = []interface{}{"Mick", "Jagger"}

	s := scenario.NewScenario()

	str, err := s.GetList(m, "rock", nil)

	if str[0] != "Mick" {
		t.Errorf("Bad result. Expected %s but was %s", "Mick", str)
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}

func TestGetList2(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = []interface{}{"Mick", "Jagger"}

	def := []interface{}{"Bruce", "Springsteen"}

	s := scenario.NewScenario()

	str, err := s.GetList(m, "test", def)

	if str[1] != "Springsteen" {
		t.Errorf("Bad result. Expected %s but was %s", "Springsteen", str)
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}

func TestGetList3(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = []interface{}{"Mick", "Jagger"}

	s := scenario.NewScenario()

	_, err := s.GetList(m, "test", nil)

	if err == nil {
		t.Errorf("Should return an error")
	} else {
		t.Logf("Error returned: %s", err.Error())
	}
}

func TestGetList4(t *testing.T) {

	def := []interface{}{"Bruce", "Springsteen"}

	s := scenario.NewScenario()

	str, err := s.GetList(nil, "test", def)

	if str[0] != "Bruce" {
		t.Errorf("Bad result. Expected %s but was %s", "Bruce", str)
	}

	if err != nil {
		t.Errorf("Should not return an error: %s", err.Error())
	}

}

func TestGetList5(t *testing.T) {

	s := scenario.NewScenario()

	_, err := s.GetList(nil, "test", nil)

	if err == nil {
		t.Errorf("Should return an error")
	} else {
		t.Logf("Error returned: %s", err.Error())
	}
}

//----

func TestGetList6(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = []interface{}{"Mick", "Jagger"}

	def := []string{"Bruce", "Springsteen"}

	s := scenario.NewScenario()

	str, err := s.GetList(m, "test", def)

	if str[1] != "Springsteen" {
		t.Errorf("Bad result. Expected %s but was %s", "Springsteen", str)
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}

func TestGetList7(t *testing.T) {

	def := []string{"Bruce", "Springsteen"}

	s := scenario.NewScenario()

	str, err := s.GetList(nil, "test", def)

	if str[0] != "Bruce" {
		t.Errorf("Bad result. Expected %s but was %s", "Bruce", str)
	}

	if err != nil {
		t.Errorf("Should not return an error: %s", err.Error())
	}

}

func TestGetList8(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = []interface{}{"Mick", "Jagger"}

	def := "Springsteen"

	s := scenario.NewScenario()

	str, err := s.GetList(m, "test", def)

	if str[0] != "Springsteen" {
		t.Errorf("Bad result. Expected %s but was %s", "Springsteen", str)
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}

func TestGetList9(t *testing.T) {

	def := "Bruce"

	s := scenario.NewScenario()

	str, err := s.GetList(nil, "test", def)

	if str[0] != "Bruce" {
		t.Errorf("Bad result. Expected %s but was %s", "Bruce", str)
	}

	if err != nil {
		t.Errorf("Should not return an error: %s", err.Error())
	}

}

func TestGetList10(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = []interface{}{"Mick", "Jagger"}

	def := 100

	s := scenario.NewScenario()

	str, err := s.GetList(m, "test", def)

	if str[0] != "100" {
		t.Errorf("Bad result. Expected %d but was %s", 100, str[0])
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}

func TestGetList11(t *testing.T) {
	m := make(map[string]interface{})
	m["rock"] = "Jagger"

	s := scenario.NewScenario()

	str, err := s.GetList(m, "rock", nil)

	if str[0] != "Jagger" {
		t.Errorf("Bad result. Expected %s but was %s", "Jagger", str)
	}

	if err != nil {
		t.Errorf("Should not return en error: %s", err.Error())
	}
}
