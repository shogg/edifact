package edifact

import "strings"

var releaseMetaChars = strings.NewReplacer(
	":", "?:",
	"+", "?+",
	"?", "??",
	"'", "?'")

// Release add '?' in front of meta characters ':', '+', '\'' to release them.
func Release(s string) string {
	return releaseMetaChars.Replace(s)
}
