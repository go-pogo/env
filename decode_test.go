// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDecoder(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilLookupper, func() {
			_ = NewDecoder(nil)
		})
	})
	t.Run("nils", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilLookupper, func() {
			_ = NewDecoder(nil, nil)
		})
	})
}

func TestDecoder_Decode(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type subj struct {
			Foo string
			Ok  bool
		}
		const input = "FOO=bar\nOK=true"

		var have subj
		dec := NewReaderDecoder(strings.NewReader(input))
		assert.NoError(t, dec.Decode(&have))
		assert.Exactly(t, subj{Foo: "bar", Ok: true}, have)
	})

	t.Run("strict", func(t *testing.T) {
		type subj struct {
			Foo string
			Ok  bool
		}
		const input = "FOO=bar\nOK=true"

		var have subj
		dec := NewReaderDecoder(strings.NewReader(input)).Strict()
		assert.NoError(t, dec.Decode(&have))
		assert.Exactly(t, subj{}, have)
	})

	t.Run("basic with ignored field", func(t *testing.T) {
		type subj struct {
			Foo    string
			Ignore bool `env:"-"`
		}
		const input = "FOO=bar\nIGNORE=true"

		var have subj
		dec := NewReaderDecoder(strings.NewReader(input))
		assert.NoError(t, dec.Decode(&have))
		assert.Exactly(t, have, subj{Foo: "bar"})
	})

	t.Run("strict with ignored field", func(t *testing.T) {
		type subj struct {
			Foo    string `env:"CUSTOM_NAME"`
			Ignore bool   `env:"-"`
		}
		const input = "CUSTOM_NAME=bar\nIGNORE=true"

		var have subj
		dec := NewReaderDecoder(strings.NewReader(input)).Strict()
		assert.NoError(t, dec.Decode(&have))
		assert.Exactly(t, subj{Foo: "bar"}, have)
	})

	t.Run("nested", func(t *testing.T) {
		type nested struct {
			Foo string
		}
		type root struct {
			Qux    string
			Nested nested
		}

		const input = "QUX=x00\nNESTED_FOO=bar"
		want := root{
			Qux:    "x00",
			Nested: nested{Foo: "bar"},
		}

		var have root
		dec := NewReaderDecoder(strings.NewReader(input))
		assert.NoError(t, dec.Decode(&have))
		assert.Exactly(t, want, have)
	})

	t.Run("inline", func(t *testing.T) {
		type inline struct {
			Foo string
		}
		type root struct {
			Empty  string
			Inline inline `env:",inline"`
		}

		const input = "FOO=bar\nINLINE_FOO=not good!"

		t.Run("loose", func(t *testing.T) {
			want := root{
				Inline: inline{Foo: "bar"},
			}

			var have root
			dec := NewReaderDecoder(strings.NewReader(input))
			assert.NoError(t, dec.Decode(&have))
			assert.Exactly(t, want, have)
		})

		t.Run("strict", func(t *testing.T) {
			var have root
			dec := NewReaderDecoder(strings.NewReader(input)).Strict()
			assert.NoError(t, dec.Decode(&have))
			assert.Exactly(t, root{}, have)
		})
	})

	t.Run("include", func(t *testing.T) {
		type included struct {
			Foo string
		}
		type root struct {
			IgnoreMe bool     `env:"-"`
			Included included `env:",include"`
		}

		const input = "IGNORE_ME=true\nINCLUDED_FOO=bar"
		want := root{
			Included: included{
				Foo: "bar",
			},
		}

		var have root
		dec := NewReaderDecoder(strings.NewReader(input))
		assert.NoError(t, dec.Decode(&have))
		assert.Exactly(t, want, have)
	})

	t.Run("nested", func(t *testing.T) {
		type deepNested struct {
			Foo string `env:"CUSTOM_NAME"`
			Bar string
		}
		type nested struct {
			Qux        string
			DeepNested deepNested `env:"DEEP"`
		}
		type root struct {
			Nested nested `env:"NESTED"`
		}

		const input = "CUSTOM_NAME=foo\nNESTED_QUX=xoo"
		want := root{
			Nested: nested{
				Qux: "xoo",
				DeepNested: deepNested{
					Foo: "foo",
				},
			},
		}

		var have root
		dec := NewReaderDecoder(strings.NewReader(input))
		assert.NoError(t, dec.Decode(&have))
		assert.Exactly(t, want, have)
	})

	t.Run("nil", func(t *testing.T) {
		assert.ErrorIs(t,
			NewDecoder(System()).Decode(nil),
			ErrStructPointerExpected,
		)
	})

	t.Run("non-pointer", func(t *testing.T) {
		type subj struct {
			Foo    string
			Ignore bool
		}
		assert.ErrorIs(t,
			NewDecoder(System()).Decode(subj{}),
			ErrStructPointerExpected,
		)
	})

	t.Run("non-struct pointer", func(t *testing.T) {
		var v int
		assert.ErrorIs(t,
			NewDecoder(System()).Decode(&v),
			ErrStructPointerExpected,
		)
	})

	t.Run("nil lookupper", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilLookupper, func() {
			var dec Decoder
			var dest struct {
				Foo    string
				Ignore bool
			}
			_ = dec.Decode(&dest)
		})
	})
}
