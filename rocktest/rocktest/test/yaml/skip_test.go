package yamlTest

import (
	"testing"
)

func TestSkip1(t *testing.T) {
	shouldFailWithMessage(t, "skip/skip.yaml", "I feel good")
}
