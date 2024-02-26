// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/env/envtag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewEncoder(t *testing.T) {
	t.Run("nil writer", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilWriter, func() { NewEncoder(nil) })
	})
}

func TestEncoder_WithWriter(t *testing.T) {
	t.Run("nil writer", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilWriter, func() {
			NewEncoder(&strings.Builder{}).WithWriter(nil)
		})
	})
}

func TestEncoder(t *testing.T) {
	urlInput, err := url.Parse("https://example.com")
	require.NoError(t, err)

	tests := map[string]struct {
		setup func(enc *Encoder)
		input any
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

		// basic types as fields
		"bool": {
			setup: takeValues,
			input: struct{ Any bool }{Any: true},
			want:  []string{"ANY=true"},
		},
		"int": {
			setup: takeValues,
			input: struct{ Any int }{Any: -123},
			want:  []string{"ANY=-123"},
		},
		"uint": {
			setup: takeValues,
			input: struct{ Any uint }{Any: 123},
			want:  []string{"ANY=123"},
		},
		"float": {
			setup: takeValues,
			input: struct{ Any float64 }{Any: math.Pi},
			want:  []string{"ANY=" + strconv.FormatFloat(math.Pi, 'g', -1, 64)},
		},
		"string": {
			setup: takeValues,
			input: struct{ Any string }{Any: "foo"},
			want:  []string{"ANY=foo"},
		},
		"url": {
			setup: takeValues,
			input: struct{ Any *url.URL }{Any: urlInput},
			want:  []string{"ANY=https://example.com"},
		},
		"duration": {
			setup: takeValues,
			input: struct{ Any time.Duration }{Any: time.Second + (time.Minute * 3)},
			want:  []string{"ANY=3m1s"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var buf strings.Builder
			enc := NewEncoder(&buf)
			if tc.setup != nil {
				tc.setup(enc)
			}

			assert.NoError(t, enc.Encode(tc.input))
			assertSimilarOutput(t, buf.String(), tc.want)
		})
	}

	t.Run("struct", func(t *testing.T) {
		type subj struct {
			Foo        string `default:"bar"`
			Flag       bool   `default:"true"`
			unexported bool   //nolint:golint,unused
		}

		t.Run("default", func(t *testing.T) {
			var buf strings.Builder
			enc := NewEncoder(&buf)

			assert.NoError(t, enc.Encode(subj{Foo: "foobar"}))
			assertSimilarOutput(t, buf.String(), []string{
				`FOO=bar`,
				`FLAG=true`,
			})
		})
		t.Run("take values", func(t *testing.T) {
			var buf strings.Builder
			enc := NewEncoder(&buf)
			enc.TakeValues = true

			assert.NoError(t, enc.Encode(subj{Foo: "foobar"}))
			assertSimilarOutput(t, buf.String(), []string{
				`FOO=foobar`,
				`FLAG=true`,
			})
		})
	})

	t.Run("nested struct", func(t *testing.T) {
		type nested struct {
			Foo        string `env:"FOO" default:"bar"`
			unexported bool   `env:"NOPE"` //nolint:golint,unused
		}
		type subj struct {
			Qux    string
			Nested nested
		}

		var buf strings.Builder
		assert.NoError(t, NewEncoder(&buf).Encode(subj{Qux: "x00"}))
		assertSimilarOutput(t, buf.String(), []string{
			`QUX=`,
			`FOO=bar`,
		})
	})
}

func assertSimilarOutput(t *testing.T, have string, want []string) {
	assert.Len(t, have, 1+len(strings.Join(want, "\n")))
	for _, line := range want {
		assert.Contains(t, have, line)
	}
}

func takeValues(enc *Encoder) {
	enc.TakeValues = true
}
