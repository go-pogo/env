// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"testing"

	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
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
				haveKey, haveVal, haveErr := Parse(input)
				assert.Empty(t, haveKey)
				assert.Empty(t, haveVal)

				if wantErr == nil {
					assert.Nil(t, haveErr)
				} else {
					assert.True(t, errors.Is(haveErr, wantErr))
				}
			})
		}
	})

	t.Run("keys", func(t *testing.T) {
		tests := map[string]struct {
			wantErr error
			wantKey string
			input   string
		}{
			"empty": {
				input:   "",
				wantErr: ErrEmptyKey,
			},
			"plain": {
				input:   "foo",
				wantKey: "foo",
			},
			"whitespace": {
				input:   "\tyolo  ",
				wantKey: "yolo",
			},
			"whitespace left": {
				input:   " bar",
				wantKey: "bar",
			},
			"whitespace right": {
				input:   "qux\t",
				wantKey: "qux",
			},
			"export": {
				input:   "export ",
				wantErr: ErrEmptyKey,
			},
			"export as valid key": {
				input:   "export",
				wantKey: "export",
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				haveKey, _, haveErr := Parse(tc.input + "=some value")
				assert.Exactly(t, tc.wantKey, haveKey)

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
			wantErr error
			wantVal Value
			input   []string // value parts
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
				wantVal: "this is 'a quote'!",
			},
			"double quote in value": {
				input: []string{
					`'"double" quotes FTW'`,
					`"\"double\" quotes FTW"`,
				},
				wantVal: `"double" quotes FTW`,
			},
			"escape sequence": {
				input:   []string{`'\\\''`},
				wantVal: `\'`,
			},
			"escape sequence 2": {
				input:   []string{`'\\\\'`},
				wantVal: `\\`,
			},
			"escape sequence 3": {
				input:   []string{`'\\\\\\'`},
				wantVal: `\\\`,
			},
			"comment at end": {
				input: []string{
					"bar # comment",
					"'bar' #comment",
					`"bar"#comment`,
				},
				wantVal: "bar",
			},
			"hash in value": {
				input: []string{
					"'#xoo' ",
					`"#xoo"`,
				},
				wantVal: "#xoo",
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
					_, haveVal, haveErr := Parse(input)

					assert.Exactly(t, tc.wantVal, haveVal, "parsing `"+input+"` failed")
					if tc.wantErr == nil {
						assert.Nil(t, haveErr)
					} else {
						assert.ErrorIs(t, haveErr, tc.wantErr)
					}
				})
			}
		}
	})
}
