// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	type fixtureBasic struct {
		Foo    string
		Ignore bool
	}

	type fixtureBasic2 struct {
		Foo    string
		Ignore bool `env:"-"`
	}

	type fixtureBasic3 struct {
		Foo    string `env:"CUSTOM_NAME,noprefix"`
		Ignore bool   `env:"-"`
	}

	type fixtureNested struct {
		Qux    string
		Nested fixtureBasic
	}

	type fixtureNested2 struct {
		Qux    string
		Nested fixtureBasic3
	}

	type fixtureInline struct {
		Qux    string
		Inline fixtureBasic2 `env:",inline"`
	}

	type fixtureNoPrefix struct {
		DeepNested fixtureNested2 `env:"DEEPNESTED"`
	}

	tests := map[string]struct {
		dec     Decoder
		input   string
		want    interface{}
		wantErr error
	}{
		"basic": {
			input:   "FOO=bar\nIGNORE=true",
			want:    &fixtureBasic{Foo: "bar", Ignore: true},
			wantErr: nil,
		},
		"basic with ignored field": {
			input:   "FOO=bar\nIGNORE=true",
			want:    &fixtureBasic2{Foo: "bar"},
			wantErr: nil,
		},
		"basic with TagsOnly": {
			dec:     Decoder{ReplaceVars: true},
			input:   "CUSTOM_NAME=bar\nIGNORE=true",
			want:    &fixtureBasic3{Foo: "bar"},
			wantErr: nil,
		},
		"nested": {
			input: "QUX=x00\nNESTED_FOO=bar",
			want:  &fixtureNested{Qux: "x00", Nested: fixtureBasic{Foo: "bar"}},
		},
		"inline": {
			input: `
QUX=x00
FOO=bar
IGNORE=true
INLINE_FOO=not used`,
			want: &fixtureInline{
				Qux: "x00",
				Inline: fixtureBasic2{
					Foo: "bar",
				},
			},
		},
		"noPrefix": {
			input: "CUSTOM_NAME=bar\nDEEPNESTED_QUX=xoo",
			want: &fixtureNoPrefix{
				DeepNested: fixtureNested2{
					Qux: "xoo",
					Nested: fixtureBasic3{
						Foo: "bar",
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := reflect.New(reflect.ValueOf(tc.want).Elem().Type()).Interface()

			tc.dec.WithLookupper(NewReader(strings.NewReader(tc.input)))
			haveErr := tc.dec.Decode(have)

			if tc.wantErr == nil {
				assert.NoError(t, haveErr)
			} else {
				assert.ErrorIs(t, tc.wantErr, haveErr)
			}
			assert.Exactly(t, tc.want, have)
		})
	}

	t.Run("nil", func(t *testing.T) {
		assert.ErrorIs(t, NewDecoder(EnvironLookup()).Decode(nil), ErrStructPointerExpected)
	})
	t.Run("non-pointer", func(t *testing.T) {
		assert.ErrorIs(t, NewDecoder(EnvironLookup()).Decode(fixtureBasic{}), ErrStructPointerExpected)
	})
	t.Run("non-struct pointer", func(t *testing.T) {
		var v int
		assert.ErrorIs(t, NewDecoder(EnvironLookup()).Decode(&v), ErrStructPointerExpected)
	})
}
