package edifact_test

import (
	"testing"

	"github.com/shogg/edifact"
	"github.com/stretchr/testify/assert"
)

func TestRelease(t *testing.T) {

	tests := []struct {
		str, want string
	}{
		{"", ""},
		{" ", " "},
		{"a", "a"},
		{"?", "??"},
		{"+", "?+"},
		{":", "?:"},
		{"'", "?'"},
		{"++", "?+?+"},
		{"Who's this?", "Who?'s this??"},
		{":+?'", "?:?+???'"},
	}

	for i := range tests {
		got := edifact.Release(tests[i].str)
		assert.Equal(t, tests[i].want, got, "test %d failed", i)
	}
}
