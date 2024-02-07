// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewReader(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilReader, func() { NewReader(nil) })
	})
}

func TestNewFileReader(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilReader, func() { NewFileReader(nil) })
	})
}
