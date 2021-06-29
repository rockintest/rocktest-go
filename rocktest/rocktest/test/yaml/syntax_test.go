package yamlTest

import (
	"testing"
)

func TestSyntaxBadStep(t *testing.T) {

	shouldFail(t, "syntax/badstep.yaml")

}

func TestSyntaxBadYaml(t *testing.T) {

	shouldFail(t, "syntax/badyaml.yaml")

}
