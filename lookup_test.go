// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestChain(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		want := EnvironLookup()
		have := Chain(want)
		assert.Exactly(t, reflect.ValueOf(want).Pointer(), reflect.ValueOf(have).Pointer())
	})
	t.Run("chain", func(t *testing.T) {
		chain1 := Chain(EnvironLookup(), EnvironLookup())
		assert.Len(t, chain1, 2)
		chain2 := Chain(chain1, EnvironLookup())
		assert.Len(t, chain2, 3)
	})
}
