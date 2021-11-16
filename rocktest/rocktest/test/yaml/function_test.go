package yamlTest

import (
	"testing"
)

func TestFunction1(t *testing.T) {
	shouldPass(t, "function/display.yaml")
}

func TestFunction2(t *testing.T) {
	shouldPass(t, "function/local.yaml")
}

func TestFunctionFail(t *testing.T) {
	shouldFail(t, "function/localFail.yaml")
}

func TestFunction3(t *testing.T) {
	shouldPass(t, "function/libcaller.yaml")
}

func TestFunctionFail2(t *testing.T) {
	shouldFail(t, "function/libcallerFail.yaml")
}

func TestFunctionNoExist(t *testing.T) {
	shouldFail(t, "function/noExistLocal.yaml")
}

func TestFunctionNoExistRemote(t *testing.T) {
	shouldFail(t, "function/noExistRemote.yaml")
}

func TestModuleNoExist(t *testing.T) {
	shouldFail(t, "function/noExistModule.yaml")
}

func TestFunction4(t *testing.T) {
	shouldPass(t, "function/submodules.yaml")
}
