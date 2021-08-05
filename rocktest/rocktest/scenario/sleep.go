package scenario

import (
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func (module *Module) Sleep(params map[string]interface{}, scenario *Scenario) error {

	paramsEx, err := scenario.ExpandMap(params)
	if err != nil {
		return err
	}

	val, err := scenario.GetString(paramsEx, "value", nil)

	if err != nil {
		return err
	}

	if val == "forever" {

		log.Debugf("Sleep forever")
		time.Sleep(time.Duration(1<<63 - 1))

	} else {

		delay, err := strconv.Atoi(val)

		if err != nil {
			return err
		}

		log.Debugf("Sleep for %d seconds", delay)
		time.Sleep(time.Duration(delay) * time.Second)

	}

	return nil
}
