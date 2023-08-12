// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"github.com/go-pogo/env/envtag"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEncoder(t *testing.T) {
	type fixtureBasic struct {
		Foo        string `env:"FOO" default:"bar"`
		unexported bool   `env:"NOPE"`
	}

	type fixtureNested struct {
		Qux    string
		Nested fixtureBasic
	}

	tests := map[string]struct {
		input interface{}
		want  []string
	}{
		"Map": {
			input: Map{
				`foo`: `bar`,
				`qux`: ``,
			},
			want: []string{
				`foo=bar`,
				`qux=`,
			},
		},
		"map": {
			input: map[string]Value{
				`foo`: `${bar}`,
				`qux`: `"xoo"`,
			},
			want: []string{
				`foo=${bar}`,
				`qux='"xoo"'`,
			},
		},
		"NamedValues": {
			input: []NamedValue{
				{Name: `foo`, Value: `12.3`},
				{Name: `qux`, Value: `'xoo'`},
			},
			want: []string{
				`foo=12.3`,
				`qux="'xoo'"`,
			},
		},
		"Tag": {
			input: []envtag.Tag{
				{Name: `FOO`, Default: `bar'n "boos"`},
			},
			want: []string{
				`FOO="bar'n \"boos\""`,
			},
		},
		"struct": {
			input: fixtureBasic{Foo: "foobar"},
			want: []string{
				`FOO=bar`,
			},
		},
		"nested struct": {
			input: fixtureNested{Qux: "x00"},
			want: []string{
				`QUX=`,
				`NESTED_FOO=bar`,
			},
		},
	}

	for name, tc := range tests {
		t.Run("default", func(t *testing.T) {
			t.Run(name, func(t *testing.T) {
				var buf bytes.Buffer
				enc := NewEncoder(&buf)

				assert.NoError(t, enc.Encode(tc.input))
				assertSimilarOutput(t, buf.String(), tc.want, "")
			})
		})
		t.Run("ExportPrefix", func(t *testing.T) {
			t.Run(name, func(t *testing.T) {
				var buf bytes.Buffer
				enc := NewEncoder(&buf)
				enc.ExportPrefix = true

				assert.NoError(t, enc.Encode(tc.input))
				assertSimilarOutput(t, buf.String(), tc.want, "export ")
			})
		})
	}

	//t.Run("TakeValues", func(t *testing.T) {
	//	if tc.wantTakeValues == nil {
	//		tc.wantTakeValues = tc.want
	//	}
	//
	//	var buf bytes.Buffer
	//	enc := NewEncoder(&buf)
	//	enc.TakeValues = true
	//
	//	assert.NoError(t, enc.Encode(tc.input))
	//	assertSimilarOutput(t, buf.String(), tc.wantTakeValues, "")
	//})
}

func assertSimilarOutput(t *testing.T, have string, want []string, prefix string) {
	assert.Len(t, have, 1+len(strings.Join(want, "\n"))+(len(prefix)*len(want)))
	for _, line := range want {
		assert.Contains(t, have, prefix+line)
	}
}
