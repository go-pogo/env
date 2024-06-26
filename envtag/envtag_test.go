// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag_IsEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.True(t, Tag{}.IsEmpty())
	})
	t.Run("not empty", func(t *testing.T) {
		assert.False(t, Tag{Name: "foo"}.IsEmpty())
	})
}

func TestTag_ShouldIgnore(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		assert.True(t, Tag{Name: ""}.ShouldIgnore())
	})
	t.Run("ignore", func(t *testing.T) {
		assert.True(t, Tag{Ignore: true}.ShouldIgnore())
	})
	t.Run("name", func(t *testing.T) {
		assert.False(t, Tag{Name: "foobar"}.ShouldIgnore())
	})
}
