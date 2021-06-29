package scenario

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func (module *Module) Assert_equals(params map[string]interface{}, scenario *Scenario) error {

	if !log.IsLevelEnabled(log.DebugLevel) {
		log.Infof("   expected= %s, actual= %s", params["expected"], params["actual"])
	}

	paramsEx := scenario.ExpandMap(params)

	expected, err := scenario.GetString(paramsEx, "expected", nil)
	if err != nil {
		return err
	}

	actual, err := scenario.GetString(paramsEx, "actual", nil)
	if err != nil {
		return err
	}

	msg, err := scenario.GetString(paramsEx, "message", "")
	if err != nil {
		return err
	}

	if expected != actual {
		errmsg := fmt.Sprintf("%s - expected %s, but was %s", msg, expected, actual)
		return errors.New(errmsg)
	}

	log.Debug("Assert OK")

	return nil
}

func (module *Module) Assert_match(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)

	expected, err := scenario.GetString(paramsEx, "expected", nil)
	if err != nil {
		return err
	}

	actual, err := scenario.GetString(paramsEx, "actual", nil)
	if err != nil {
		return err
	}

	msg, err := scenario.GetString(paramsEx, "message", "")
	if err != nil {
		return err
	}

	matched, err := regexp.Match(expected, []byte(actual))
	if err != nil {
		return err
	}

	if !matched {
		errmsg := fmt.Sprintf("%s - regex %s does not match %s", msg, expected, actual)
		return errors.New(errmsg)
	}

	return nil
}

func (module *Module) Assert_fail(params map[string]interface{}, scenario *Scenario) error {

	paramsEx := scenario.ExpandMap(params)
	ret, err := scenario.GetString(paramsEx, "value", "")

	if err != nil {
		return err
	}

	return errors.New(ret)

}
