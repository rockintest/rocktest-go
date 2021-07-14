package yamlTest

import (
	"testing"
)

func TestMock1(t *testing.T) {
	shouldPass(t, "http/mock.yaml")
}

func TestMock2(t *testing.T) {
	shouldFail(t, "http/mockFail1.yaml")
}

func TestMock3(t *testing.T) {
	shouldFail(t, "http/mockFail2.yaml")
}

func TestMock4(t *testing.T) {
	shouldFail(t, "http/mockFail3.yaml")
}

func TestMock5(t *testing.T) {
	shouldPass(t, "http/mock5.yaml")
}

func TestMock6(t *testing.T) {
	shouldPass(t, "http/mock6.yaml")
}
