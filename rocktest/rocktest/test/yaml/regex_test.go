package yamlTest

import (
	"testing"
)

func TestRegex1(t *testing.T) {
	shouldPass(t, "regex/regexOK.yaml")
}

func TestRegex2(t *testing.T) {
	shouldPass(t, "regex/regexCompact.yaml")
}

func TestRegex3(t *testing.T) {
	shouldPass(t, "regex/regexGroup.yaml")
}

func TestRegex4(t *testing.T) {
	shouldPass(t, "regex/regexMultiline.yaml")
}
