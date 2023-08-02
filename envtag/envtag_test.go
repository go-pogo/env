// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTag_IsEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.True(t, Tag{}.IsEmpty())
	})
	t.Run("not empty", func(t *testing.T) {
		assert.False(t, Tag{Name: "foo"}.IsEmpty())
	})
}
