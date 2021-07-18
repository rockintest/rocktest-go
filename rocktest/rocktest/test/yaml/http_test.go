package yamlTest

import (
	"testing"
)

func TestHttpExpect1(t *testing.T) {
	shouldPass(t, "http/httpExpect1.yaml")
}

func TestHttpPrefix(t *testing.T) {
	shouldPass(t, "http/httpPrefix.yaml")
}

func TestHttpExpect2(t *testing.T) {
	shouldPass(t, "http/httpExpect2.yaml")
}

func TestHttpExpectFail1(t *testing.T) {
	shouldFail(t, "http/httpExpectFail1.yaml")
}

func TestHttpExpectFail2(t *testing.T) {
	shouldFail(t, "http/httpExpectFail2.yaml")
}

func TestHttpMethods(t *testing.T) {
	shouldPass(t, "http/httpMethods.yaml")
}
