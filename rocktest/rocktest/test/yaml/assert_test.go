package yamlTest

import (
	"testing"
)

func TestAssert1(t *testing.T) {
	shouldPass(t, "assert/assert.yaml")
}

func TestAssert2(t *testing.T) {
	shouldFail(t, "assert/assert-fail1.yaml")
}

func TestAssert3(t *testing.T) {
	shouldFail(t, "assert/assert-fail2.yaml")
}

func TestAssert4(t *testing.T) {
	shouldFail(t, "assert/assert-syntax2.yaml")
}

func TestAssert5(t *testing.T) {
	shouldFail(t, "assert/assert-syntax3.yaml")
}

func TestAssert6(t *testing.T) {
	shouldPass(t, "assert/assert-regex.yaml")
}

func TestAssert7(t *testing.T) {
	shouldFail(t, "assert/assert-regex-fail.yaml")
}

func TestAssert8(t *testing.T) {
	shouldPass(t, "assert/assert-set1.yaml")
}

func TestAssert9(t *testing.T) {
	shouldFail(t, "assert/assert-set2.yaml")
}
