package yamlTest

import (
	"testing"
)

func TestSleepKO(t *testing.T) {

	shouldFail(t, "sleep/sleepKO.yaml")

}

func TestSleepOK(t *testing.T) {

	shouldPass(t, "sleep/sleepOK.yaml")

}
