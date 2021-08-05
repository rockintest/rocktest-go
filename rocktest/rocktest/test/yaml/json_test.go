package yamlTest

import (
	"testing"
)

func TestJson1(t *testing.T) {
	shouldPass(t, "json/json.yaml")
}

func TestJson2(t *testing.T) {
	shouldFail(t, "json/badjson.yaml")
}

func TestJson3(t *testing.T) {
	shouldFail(t, "json/missingParam.yaml")
}

func TestJson4(t *testing.T) {
	shouldFail(t, "json/badjsonInline.yaml")
}
