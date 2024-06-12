// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileEncoder(t *testing.T) {
	t.Run("nil file", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilFile, func() { NewEncoder(nil) })
	})
}
