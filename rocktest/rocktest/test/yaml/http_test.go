package yamlTest

import (
	"testing"
)

func TestHttpPrefix(t *testing.T) {
	shouldPass(t, "http/httpPrefix.yaml")
}

func TestHttpExpect1(t *testing.T) {
	shouldPass(t, "http/httpExpect1.yaml")
}

func TestHttpExpectFail1(t *testing.T) {
	shouldFail(t, "http/httpExpectFail1.yaml")
}
