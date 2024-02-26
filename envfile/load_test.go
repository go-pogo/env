// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"github.com/go-pogo/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"testing/fstest"
)

func TestLoadFS(t *testing.T) {
	t.Run("nil fsys", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilFsys, func() {
			_ = LoadFS(nil, "filename")
		})
	})
	t.Run("file does not exist", func(t *testing.T) {
		assert.Error(t, LoadFS(fstest.MapFS{}, "filename"))
	})

	t.Run("load", func(t *testing.T) {
		have := env.Map{"QUX": "x00"}
		env.Use(have)
		defer env.Use(env.System())

		fsys := fstest.MapFS{"x": {Data: []byte("FOO=bar\nQUX=nop")}}
		assert.NoError(t, LoadFS(fsys, "x"))
		assert.Equal(t, env.Map{"FOO": "bar", "QUX": "x00"}, have)
	})
}

func TestOverloadFS(t *testing.T) {
	t.Run("nil fsys", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilFsys, func() {
			_ = OverloadFS(nil, "filename")
		})
	})
	t.Run("file does not exist", func(t *testing.T) {
		assert.Error(t, OverloadFS(fstest.MapFS{}, "filename"))
	})

	t.Run("overload", func(t *testing.T) {
		have := env.Map{"QUX": "x00"}
		env.Use(have)
		defer env.Use(env.System())

		fsys := fstest.MapFS{"y": {Data: []byte("FOO=bar\nQUX=overload")}}
		assert.NoError(t, OverloadFS(fsys, "y"))
		assert.Equal(t, env.Map{"FOO": "bar", "QUX": "overload"}, have)
	})
}
