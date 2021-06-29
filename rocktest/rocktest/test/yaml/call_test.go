package yamlTest

import (
	"testing"
)

func TestCall1(t *testing.T) {
	shouldPass(t, "call/libtest.yaml")
}

func TestCall2(t *testing.T) {
	shouldPass(t, "call/vartest.yaml")
}

func TestCall3(t *testing.T) {
	shouldPass(t, "call/vartestParams.yaml")
}

func TestCall4(t *testing.T) {
	shouldPass(t, "call/libtestReturn.yaml")
}

func TestCall5(t *testing.T) {
	shouldPass(t, "call/context.yaml")
}

func TestCall6(t *testing.T) {
	shouldPass(t, "call/libchecktest.yaml")
}

func TestCall7(t *testing.T) {
	shouldFail(t, "call/libchecktestFail.yaml")
}
