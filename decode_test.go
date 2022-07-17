// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleUnmarshal() {
	type Envs struct {
		Foo,
		Bar string
		Qux struct {
			Xoo string
		}
	}

	var data = `
FOO=bar
 BAR='baz'

# ignore me
QUX_XOO="#xoo xoo"
COMMENT=#ignored`

	var envs Envs
	if err := Unmarshal([]byte(data), &envs); err != nil {
		panic(err)
	}

	spew.Dump(envs)
	// Output:
	// (env.Envs) {
	//  Foo: (string) (len=3) "bar",
	//  Bar: (string) (len=3) "baz",
	//  Qux: (struct { Xoo string }) {
	//   Xoo: (string) (len=8) "#xoo xoo"
	//  }
	// }
}

func TestUnmarshal(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type fixture struct {
			Foo    string
			Ignore bool `env:"-"`
		}

		var have fixture
		haveErr := Unmarshal([]byte("FOO=bar\nIGNORE=true"), &have)
		assert.Exactly(t, fixture{Foo: "bar"}, have)
		assert.Nil(t, haveErr)
	})
}

func TestDecoder_Lookup(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		dec := NewDecoder(new(bytes.Buffer), DefaultOptions)
		haveVal, haveOk := dec.Lookup("nope")
		assert.Equal(t, Value(""), haveVal)
		assert.False(t, haveOk)
	})
	t.Run("fallback", func(t *testing.T) {
		dec := NewDecoder(new(bytes.Buffer), DefaultOptions)
		dec.Fallback(LookupperFunc(LookupEnv))
		haveVal, haveOk := dec.Lookup("PATH")
		assert.Equal(t, Getenv("PATH"), haveVal)
		assert.True(t, haveOk)
	})
	t.Run("error", func(t *testing.T) {
		dec := NewDecoder(bytes.NewBuffer([]byte("FOO='missing quote")), DefaultOptions)
		require.Nil(t, dec.Err())
		haveVal, haveOk := dec.Lookup("FOO")
		assert.Equal(t, Value(""), haveVal)
		assert.False(t, haveOk)

		var targetErr = new(LookupError)
		assert.ErrorAs(t, dec.Err(), &targetErr)
		assert.NotNil(t, targetErr)
		assert.Equal(t, "FOO", targetErr.Key)
	})
}

func TestDecoder_Map(t *testing.T) {
	wantMap := Map{"FOO": "", "bar": "3.14"}
	input := "FOO=\nbar=3.14"

	dec := NewDecoder(bytes.NewBuffer([]byte(input)), DefaultOptions)
	haveMap, haveErr := dec.Map()
	assert.Equal(t, wantMap, haveMap)
	assert.Nil(t, haveErr)

	t.Run("should match decode", func(t *testing.T) {
		haveMap2 := make(Map)
		dec = NewDecoder(bytes.NewBuffer([]byte(input)), DefaultOptions)
		assert.Nil(t, dec.Decode(haveMap2))
		assert.Equal(t, wantMap, haveMap2)
	})
}
