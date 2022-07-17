// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	se1, se2 := os.LookupEnv("GOROOT")

	tests := map[string]struct {
		maps    []Lookupper
		key     string
		wantVal Value
		wantOk  bool
	}{
		"empty map": {
			maps:   []Lookupper{},
			key:    "foo",
			wantOk: false,
		},
		"one map": {
			maps:    []Lookupper{map1},
			key:     "foo",
			wantVal: "bar",
			wantOk:  true,
		},
		"one map, invalid key": {
			maps:   []Lookupper{map1},
			key:    "bar",
			wantOk: false,
		},
		"two maps": {
			maps:    []Lookupper{map1, map2},
			key:     "clark",
			wantVal: "superman",
			wantOk:  true,
		},
		"two maps, invalid key": {
			maps:   []Lookupper{map1, map2},
			key:    "peter",
			wantOk: false,
		},
		"system env": {
			maps:    []Lookupper{map1, map2, LookupperFunc(LookupEnv)},
			key:     "GOROOT",
			wantVal: Value(se1),
			wantOk:  se2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			haveVal, haveOk := Lookup(tc.key, tc.maps...)
			assert.Exactly(t, tc.wantVal, haveVal)
			assert.Exactly(t, tc.wantOk, haveOk)
		})
	}
}
