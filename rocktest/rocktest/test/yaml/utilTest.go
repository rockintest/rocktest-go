package yamlTest

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"io.rocktest/rocktest/scenario"
)

func initLog() {
	log.SetLevel(log.DebugLevel)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
}

func exec(t *testing.T, scen string) error {
	initLog()

	s := scenario.NewScenario()
	err := s.Run("scen/" + scen)

	t.Logf("Scenario return:\n%v", err)

	return err
}

func shouldPass(t *testing.T, scen string) {

	err := exec(t, scen)

	if err != nil {
		t.Errorf("Error unexpected")
	}
}

func shouldFail(t *testing.T, scen string) {

	err := exec(t, scen)

	if err == nil {
		t.Errorf("Error expected")
	}
}

func shouldFailWithMessage(t *testing.T, scen string, msg string) {

	err := exec(t, scen)

	if err == nil {
		t.Errorf("Error expected")
		return
	}

	errMap := make(map[string]interface{})

	yaml.Unmarshal([]byte(err.Error()), errMap)

	if errMap["error"].(string) != msg {
		t.Errorf("Bad message type. Expected '%s' but was %s", msg, errMap["error"])
	}

}

func shouldFailAtStep(t *testing.T, scen string, step int) {

	err := exec(t, scen)

	if err == nil {
		t.Errorf("Error expected")
		return
	}

	errMap := make(map[string]interface{})

	yaml.Unmarshal([]byte(err.Error()), errMap)

	if errMap["stepNumber"].(int) != step {
		t.Errorf("Bad step number. Expected %d but was %v", step, errMap["stepNumber"])
	}

}
