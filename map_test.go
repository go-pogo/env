// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestMap_Set(t *testing.T) {
	m := make(Map, 1)
	assert.NoError(t, m.Set("foo", "bar"))
	assert.Len(t, m, 1)
	assert.Equal(t, Value("bar"), m["foo"])
}

func TestMap_Get(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		assert.Equal(t, Value("bar"), Map{"foo": "bar"}.Get("foo"))
	})
	t.Run("not found", func(t *testing.T) {
		assert.Equal(t, Value(""), Map{}.Get("foo"))
	})
}

func TestMap_Has(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		assert.True(t, Map{"foo": ""}.Has("foo"))
	})
	t.Run("not found", func(t *testing.T) {
		assert.False(t, Map{}.Has("foo"))
	})
}

func TestMap_Lookup(t *testing.T) {
	m := Map{"foo": "bar"}
	t.Run("found", func(t *testing.T) {
		haveVal, haveErr := m.Lookup("foo")
		assert.NoError(t, haveErr)
		assert.Equal(t, Value("bar"), haveVal)
	})
	t.Run("not found", func(t *testing.T) {
		haveVal, haveErr := m.Lookup("bar")
		assert.ErrorIs(t, haveErr, ErrNotFound)
		assert.Equal(t, Value(""), haveVal)
	})
}

func TestMap_Merge(t *testing.T) {
	tests := map[string]struct {
		env   Map
		merge map[string]string
		want  Map
	}{
		"append": {
			env:   Map{"foo": "bar"},
			merge: map[string]string{"qux": "xoo"},
			want:  Map{"foo": "bar", "qux": "xoo"},
		},
		"replace": {
			env:   Map{"foo": "bar", "qux": "xoo"},
			merge: map[string]string{"qux": "bar", "foo": "xoo"},
			want:  Map{"foo": "xoo", "qux": "bar"},
		},
		"merge": {
			env:   Map{"foo": "bar", "bar": "baz"},
			merge: map[string]string{"baz": "foo", "bar": "qux"},
			want:  Map{"foo": "bar", "bar": "qux", "baz": "foo"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.env.Merge(tc.merge)
			assert.Exactly(t, tc.want, tc.env)

			have := make(Map, len(tc.env))
			have.MergeValues(tc.env)
			assert.Exactly(t, tc.want, have)
		})
	}
}

func TestMap_Clone(t *testing.T) {
	src := Map{"foo": "bar", "bar": "baz"}
	clone := src.Environ()

	assert.Equal(t, src, clone)
	assert.NotSame(t, src, clone)

	src["foo"] = "qux"
	assert.Equal(t, Value("bar"), clone["foo"])
}

func randKey() string {
	return "somewhat_random_key_" + strconv.FormatInt(time.Now().Unix(), 10)
}
