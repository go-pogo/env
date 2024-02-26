// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFileReader(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilFile, func() {
			NewReader(nil)
		})
	})
}

func TestOpenFS(t *testing.T) {
	t.Run("nil fsys", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilFsys, func() {
			OpenFS(nil, "")
		})
	})
}
