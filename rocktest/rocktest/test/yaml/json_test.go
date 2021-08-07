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

func TestJsonCheck1(t *testing.T) {
	shouldPass(t, "json/jsoncheck.yaml")
}

func TestJsonCheckBadSyntax(t *testing.T) {
	shouldFail(t, "json/jsoncheck-badsyntax.yaml")
}

func TestJsonCheckBadSyntax2(t *testing.T) {
	shouldFail(t, "json/jsoncheck-badsyntax2.yaml")
}

func TestJsonCheckFailEquals(t *testing.T) {
	shouldFail(t, "json/jsoncheck-failequals.yaml")
}

func TestJsonCheckFailMatch(t *testing.T) {
	shouldFail(t, "json/jsoncheck-failmatch.yaml")
}
