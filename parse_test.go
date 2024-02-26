// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		tests := map[string]error{
			"":                   nil,
			"#comment":           nil,
			" # another comment": nil,
			"export":             ErrInvalidFormat,
			"export ":            ErrInvalidFormat,
		}

		for input, wantErr := range tests {
			t.Run(input, func(t *testing.T) {
				have, haveErr := Parse(input)
				assert.Equal(t, NamedValue{}, have)

				if wantErr == nil {
					assert.Nil(t, haveErr)
				} else {
					assert.True(t, errors.Is(haveErr, wantErr))
				}
			})
		}
	})

	t.Run("names", func(t *testing.T) {
		tests := map[string]struct {
			input   string
			want    string
			wantErr error
		}{
			"empty": {
				input:   "",
				wantErr: ErrEmptyKey,
			},
			"plain": {
				input: "foo",
				want:  "foo",
			},
			"whitespace": {
				input: "\tyolo  ",
				want:  "yolo",
			},
			"whitespace left": {
				input: " bar",
				want:  "bar",
			},
			"whitespace right": {
				input: "qux\t",
				want:  "qux",
			},
			"export": {
				input:   "export ",
				wantErr: ErrEmptyKey,
			},
			"export as valid key": {
				input: "export",
				want:  "export",
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				have, haveErr := Parse(tc.input + "=some value")
				assert.Exactly(t, tc.want, have.Name)

				if tc.wantErr == nil {
					assert.Nil(t, haveErr)
				} else {
					assert.True(t, errors.Is(haveErr, tc.wantErr))
				}
			})
		}
	})

	t.Run("values", func(t *testing.T) {
		tests := map[string]struct {
			input   []string // value parts
			want    Value
			wantErr error
		}{
			"empty": {
				input: []string{
					"",
					" ",
					"\t\t",
					`''`,
					`""`,
					"#comment",
					"\t# comment",
					"'' # another comment",
				},
			},
			"single quote in value": {
				input: []string{
					"this is 'a quote'!",
					`'this is \'a quote\'!'`,
					`"this is 'a quote'!"`,
				},
				want: "this is 'a quote'!",
			},
			"double quote in value": {
				input: []string{
					`'"double" quotes FTW'`,
					`"\"double\" quotes FTW"`,
				},
				want: `"double" quotes FTW`,
			},
			"escape sequence": {
				input: []string{`'\\\''`},
				want:  `\'`,
			},
			"escape sequence 2": {
				input: []string{`'\\\\'`},
				want:  `\\`,
			},
			"escape sequence 3": {
				input: []string{`'\\\\\\'`},
				want:  `\\\`,
			},
			"comment at end": {
				input: []string{
					"bar # comment",
					"'bar' #comment",
					`"bar"#comment`,
				},
				want: "bar",
			},
			"hash in value": {
				input: []string{
					"'#xoo' ",
					`"#xoo"`,
				},
				want: "#xoo",
			},
			"missing endquote": {
				input: []string{
					`'`,
					`"`,
					`"'`,
					`'"`,
				},
				wantErr: ErrMissingEndQuote,
			},
		}

		for name, tc := range tests {
			for _, input := range tc.input {
				t.Run(name, func(t *testing.T) {
					input = "someKey=" + input
					have, haveErr := Parse(input)

					assert.Exactly(t, tc.want, have.Value, "parsing `"+input+"` failed")
					if tc.wantErr == nil {
						assert.NoError(t, haveErr)
					} else {
						assert.ErrorIs(t, haveErr, tc.wantErr)
					}
				})
			}
		}
	})
}
