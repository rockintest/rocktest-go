package text

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func initLog() {

	log.SetLevel(log.TraceLevel)
}

func AssertEquals(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf("No match, expecting %s, but was %s", expected, actual)
	}
}

func AssertTrue(t *testing.T, actual bool) {
	if !actual {
		t.Errorf("Should be true")
	}
}

func AssertFalse(t *testing.T, actual bool) {
	if actual {
		t.Errorf("Should be false")
	}
}

func TestQuoter1(t *testing.T) {
	initLog()

	var q ParamQuoter
	res, found := q.Lookup(`$module(p1,p2)`)
	expected := `${$module(<<[p1]>>,<<[p2]>>)}`

	AssertEquals(t, expected, res)
	AssertTrue(t, found)
}

func TestQuoter2(t *testing.T) {
	initLog()

	var q ParamQuoter
	_, found := q.Lookup(`$module()`)

	AssertFalse(t, found)
}

func TestQuoter3(t *testing.T) {
	initLog()

	var q ParamQuoter
	res, found := q.Lookup(`$module(p1,p2).ext`)
	expected := `${$module(<<[p1]>>,<<[p2]>>).ext}`

	AssertEquals(t, expected, res)
	AssertTrue(t, found)
}

func TestQuoter4(t *testing.T) {
	initLog()

	var q ParamQuoter
	_, found := q.Lookup(`$module(<<[p1]>>,<<[p2]>>)`)

	AssertFalse(t, found)
}

func TestQuoter5(t *testing.T) {
	initLog()

	var q ParamQuoter
	_, found := q.Lookup(`module(a,b)`)

	AssertFalse(t, found)
}

func TestQuoter6(t *testing.T) {
	initLog()

	var q ParamQuoter
	_, found := q.Lookup(`rock`)

	AssertFalse(t, found)
}
