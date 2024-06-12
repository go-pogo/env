// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReader(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilReader, func() { NewReader(nil) })
	})
}
