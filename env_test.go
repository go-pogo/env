// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnviron(t *testing.T) {
	source := os.Environ()
	target := Environ()

	wantLen := len(source)
	for _, e := range source {
		if e[0] == '=' {
			// exclude env variables like =::=:: or =R:=R:\\
			wantLen--
		}
	}

	assert.Exactly(t, wantLen, len(target))

	for k, v := range target {
		have := k + "=" + v.String()
		assert.Contains(t, source, have, "raw env string `%s` not in os.Environ()", have)
	}
}

func TestLookup(t *testing.T) {
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
			src:     []Lookupper{map1, map2, EnvironLookup()},
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
				assert.ErrorIs(t, haveErr, ErrNotFound)
			}
		})
	}
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
