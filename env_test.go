package env

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseLine(t *testing.T) {
	tests := map[string]Map{
		"=::=::": {"=::": "::"},
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			have := make(Map, 0)
			ParseLine(input, have)
			if !cmp.Equal(have, want) {
				t.Fail()
			}
		})
	}
}
