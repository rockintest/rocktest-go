package text

import (
	"bytes"
)

const (
	START = iota
	DOLLAR
	BRACKET
	DOLLARBRACKET
)

type StringSubstitutor struct {
	lookup Lookuper
	state  int
	level  int
}

func (subst *StringSubstitutor) Replace(s string) string {

	subst.state = START
	subst.level = -1

	var ret bytes.Buffer
	var varName []bytes.Buffer = make([]bytes.Buffer, 10)

	for _, r := range s {
		c := string(r)
		switch subst.state {
		case START:
			switch c {
			case "$":
				subst.state = DOLLAR
			default:
				ret.WriteString(c)
			}
		case DOLLAR:
			switch c {
			case "{":
				subst.level++
				subst.state = BRACKET
			case "$":
				ret.WriteString(c)
				subst.state = START
			default:
				ret.WriteString("$")
				ret.WriteString(c)
				subst.state = START
			}
		case DOLLARBRACKET:
			switch c {
			case "{":
				subst.level++
				subst.state = BRACKET
			case "$":
				varName[subst.level].WriteString(c)
				subst.state = BRACKET
			default:
				varName[subst.level].WriteString("$")
				varName[subst.level].WriteString(c)
				subst.state = BRACKET
			}
		case BRACKET:
			switch c {
			case "$":
				subst.state = DOLLARBRACKET
			case "}":

				varNameString := varName[subst.level].String()
				varName[subst.level].Reset()
				subst.level--

				// Last closing bracket ?
				if subst.level == -1 {
					if varValue, ok := subst.lookup.Lookup(varNameString); ok {
						ret.WriteString(varValue)
					} else {
						ret.WriteString("${" + varNameString + "}")
					}
					subst.state = START
				} else {
					if varValue, ok := subst.lookup.Lookup(varNameString); ok {
						varName[subst.level].WriteString(varValue)
					} else {
						varName[subst.level].WriteString("${" + varNameString + "}")
					}
				}
			default:
				varName[subst.level].WriteString(c)
			}
		}

	}

	return ret.String()
}

func NewStringSubstitutorByLookuper(lookup Lookuper) *StringSubstitutor {
	ret := new(StringSubstitutor)
	ret.state = START
	ret.lookup = lookup
	return ret
}

func NewStringSubstitutorByMap(themap map[string]string) *StringSubstitutor {

	mapLookup := NewMapLookup(themap)
	ret := NewStringSubstitutorByLookuper(mapLookup)

	return ret
}
