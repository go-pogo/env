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
