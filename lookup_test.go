// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestLookupFrom(t *testing.T) {
	map1 := Map{"foo": "bar", "qux": "xoo"}
	map2 := Map{"bruce": "batman", "clark": "superman"}

	var se2 error
	se1, ok := os.LookupEnv("GOROOT")
	if !ok {
		se2 = ErrNotFound
	}

	tests := map[string]struct {
		src     []Lookupper
		key     string
		wantVal Value
		wantErr error
	}{
		"empty map": {
			src:     []Lookupper{},
			key:     "foo",
			wantErr: ErrNotFound,
		},
		"one map": {
			src:     []Lookupper{map1},
			key:     "foo",
			wantVal: "bar",
		},
		"one map, invalid key": {
			src:     []Lookupper{map1},
			key:     "bar",
			wantErr: ErrNotFound,
		},
		"two maps": {
			src:     []Lookupper{map1, map2},
			key:     "clark",
			wantVal: "superman",
		},
		"two maps, invalid key": {
			src:     []Lookupper{map1, map2},
			key:     "peter",
			wantErr: ErrNotFound,
		},
		"system env": {
			src:     []Lookupper{map1, map2, System()},
			key:     "GOROOT",
			wantVal: Value(se1),
			wantErr: se2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			haveVal, haveErr := Lookup(tc.key, tc.src...)
			assert.Exactly(t, tc.wantVal, haveVal)

			if tc.wantErr == nil {
				assert.Nil(t, haveErr)
			} else {
				assert.True(t, IsNotFound(haveErr))
			}
		})
	}
}

func TestChain(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		want := System()
		have := Chain(want)
		assert.Exactly(t, reflect.ValueOf(want).Pointer(), reflect.ValueOf(have).Pointer())
	})
	t.Run("chain", func(t *testing.T) {
		chain1 := Chain(System(), System())
		assert.Len(t, chain1, 2)
		chain2 := Chain(chain1, System())
		assert.Len(t, chain2, 3)
	})
	t.Run("nil", func(t *testing.T) {
		assert.Equal(t, chainLookupper{}, Chain(nil))
	})
	t.Run("nils", func(t *testing.T) {
		assert.Equal(t, chainLookupper{}, Chain(nil, nil))
	})
}

func TestChain_Lookup(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		_, err := Chain(nil).Lookup("doesnt matter")
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
