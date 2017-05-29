package grafton

import (
	"regexp"
)

// NameRegexpString is the string form of a regular expression for determining
// if a credential name is valid or not.
//
// Specified in Shell and Utilities volume of IEEE 1003.1-2001.
const NameRegexpString = "^[A-Z][A-Z0-9_]{0,127}$"

var nameRegexp = regexp.MustCompile(NameRegexpString)

// ValidCredentialName returns true or false depending on whether or not the
// given name is a valid credential name.
func ValidCredentialName(name string) bool {
	return nameRegexp.MatchString(name)
}
