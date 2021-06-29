package yamlTest

import (
	"testing"
)

func TestVar1(t *testing.T) {
	shouldFail(t, "var/varKO.yaml")
}

func TestVar2(t *testing.T) {
	shouldPass(t, "var/var.yaml")
}

func TestVar3(t *testing.T) {
	shouldPass(t, "var/varConcat.yaml")
}

func TestVar4(t *testing.T) {
	shouldPass(t, "var/vartest2.yaml")
}

func TestVar5(t *testing.T) {
	shouldPass(t, "var/varBuiltin.yaml")
}

func TestExpr1(t *testing.T) {
	shouldPass(t, "var/varSubstTest.yaml")
}
