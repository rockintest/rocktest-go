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

func TestMock7(t *testing.T) {
	shouldFailAtStep(t, "http/mock7.yaml", 2)
}

func TestMockHeaders(t *testing.T) {
	shouldPass(t, "http/mockHeaders.yaml")
}

func TestMockHeaders2(t *testing.T) {
	shouldPass(t, "http/mockHeaders2.yaml")
}
