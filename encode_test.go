// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"github.com/go-pogo/env/envtag"
	"github.com/stretchr/testify/assert"
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
			NewEncoder(&bytes.Buffer{}).WithWriter(nil)
		})
	})
}

func TestEncoder(t *testing.T) {
	type fixtureBasic struct {
		Foo        string `env:"FOO" default:"bar"`
		unexported bool   `env:"NOPE"` //nolint:golint,unused
	}

	type fixtureNested struct {
		Qux    string
		Nested fixtureBasic
	}

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
			setup: func(enc *Encoder) {
				enc.ExportPrefix = true
			},
			input: map[string]Value{
				`foo`: `${bar}`,
				`qux`: `"xoo"`,
			},
			want: []string{
				`export foo=${bar}`,
				`export qux='"xoo"'`,
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
				`FOO=bar`,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := NewEncoder(&buf)
			if tc.setup != nil {
				tc.setup(enc)
			}

			assert.NoError(t, enc.Encode(tc.input))
			assertSimilarOutput(t, buf.String(), tc.want)
		})
	}

	typeTests := map[string]struct {
		exec func(enc *Encoder) error
		want string
	}{
		"bool": {
			exec: func(enc *Encoder) error {
				return enc.Encode(struct{ Any bool }{Any: true})
			},
			want: "true",
		},
		"int": {
			exec: func(enc *Encoder) error {
				return enc.Encode(struct{ Any int }{Any: -123})
			},
			want: "-123",
		},
		"uint": {
			exec: func(enc *Encoder) error {
				return enc.Encode(struct{ Any uint }{Any: 123})
			},
			want: "123",
		},
		"float": {
			exec: func(enc *Encoder) error {
				return enc.Encode(struct{ Any float64 }{Any: math.Pi})
			},
			want: strconv.FormatFloat(math.Pi, 'g', -1, 64),
		},
		"string": {
			exec: func(enc *Encoder) error {
				return enc.Encode(struct{ Any string }{Any: "foo"})
			},
			want: "foo",
		},
		"url": {
			exec: func(enc *Encoder) error {
				u, err := url.Parse("https://example.com")
				if err != nil {
					panic(err)
				}
				return enc.Encode(struct{ Any *url.URL }{Any: u})
			},
			want: "https://example.com",
		},
		"duration": {
			exec: func(enc *Encoder) error {
				return enc.Encode(struct{ Any time.Duration }{Any: time.Second + (time.Minute * 3)})
			},
			want: "3m1s",
		},
	}

	t.Run("TakeValues", func(t *testing.T) {
		for name, tc := range typeTests {
			t.Run(name, func(t *testing.T) {
				var buf bytes.Buffer
				enc := NewEncoder(&buf)
				enc.TakeValues = true

				assert.NoError(t, tc.exec(enc))
				assert.Equal(t, "ANY="+tc.want+"\n", buf.String())
			})
		}
	})
}

func assertSimilarOutput(t *testing.T, have string, want []string) {
	assert.Len(t, have, 1+len(strings.Join(want, "\n")))
	for _, line := range want {
		assert.Contains(t, have, line)
	}
}
