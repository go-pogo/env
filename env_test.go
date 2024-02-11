// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetenv(t *testing.T) {
	t.Run("err", func(t *testing.T) {
		wantErr := os.Setenv("", "foobar")
		haveErr := Setenv("", "foobar")
		assert.EqualError(t, haveErr, wantErr.Error())
	})
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

func TestEnvironLookup(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		_, err := EnvironLookup().Lookup(randKey())
		assert.True(t, IsNotFound(err))
	})
}
