package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceVars(t *testing.T) {
	tests := map[string]struct {
		input   Map
		want    Map
		wantErr error
	}{
		"none": {
			input: map[string]Value{"foo": "bar"},
			want:  map[string]Value{"foo": "bar"},
		},
		"basic": {
			input: map[string]Value{"foo": `$bar`, "bar": "baz"},
			want:  map[string]Value{"foo": "baz", "bar": "baz"},
		},
		"double": {
			input: map[string]Value{"foo": `$bar and $bar`, "bar": "baz"},
			want:  map[string]Value{"foo": "baz and baz", "bar": "baz"},
		},
		"missing": {
			input: map[string]Value{"foo": `$missing var`, "bar": "baz"},
			want:  map[string]Value{"foo": "$missing var", "bar": "baz"},
		},
		"bash style": {
			input: map[string]Value{"foo": `${bar}`, "bar": "baz"},
			want:  map[string]Value{"foo": "baz", "bar": "baz"},
		},
		"bash style with default": {
			input: map[string]Value{"foo": `look at ${baz:-that}`, "bar": "baz"},
			want:  map[string]Value{"foo": "look at that", "bar": "baz"},
		},
		"dependencies": {
			input: map[string]Value{"foo": `$bar`, "bar": `${qux:-baz}`},
			want:  map[string]Value{"foo": "baz", "bar": "baz"},
		},
		"multiple dependencies": {
			input: map[string]Value{"foo": `some $bar ${qux} thing`, "bar": `$qux`, "qux": "xoo"},
			want:  map[string]Value{"foo": "some xoo xoo thing", "bar": "xoo", "qux": "xoo"},
		},
		"circular dependency": {
			input:   map[string]Value{"foo": `$bar`, "bar": `$foo`},
			want:    map[string]Value{"foo": `$bar`, "bar": `$foo`},
			wantErr: ErrCircularDependency,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have, haveErr := ReplaceAll(tc.input)
			assert.Equal(t, tc.want, have)
			if tc.wantErr == nil {
				assert.NoError(t, haveErr)
			} else {
				assert.ErrorIs(t, haveErr, tc.wantErr)
			}
		})
	}
}
