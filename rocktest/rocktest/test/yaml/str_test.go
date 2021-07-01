package yamlTest

import (
	"testing"
)

func TestStr1(t *testing.T) {
	shouldPass(t, "str/toupper.yaml")
}

func TestStr2(t *testing.T) {
	shouldPass(t, "str/tolower.yaml")
}

func TestStr3(t *testing.T) {
	shouldPass(t, "str/tolowerInline.yaml")
}

func TestStr4(t *testing.T) {
	shouldPass(t, "str/toupperInline.yaml")
}
