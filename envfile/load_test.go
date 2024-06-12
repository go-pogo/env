// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"testing"
	"testing/fstest"

	"github.com/go-pogo/env"
	"github.com/go-pogo/env/envtest"
	"github.com/stretchr/testify/assert"
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
		envs := envtest.Prepare(env.Map{"QUX": "x00"})
		defer envs.Restore()

		fsys := fstest.MapFS{"x": {Data: []byte("FOO=bar\nQUX=nop")}}
		assert.NoError(t, LoadFS(fsys, "x"))
		assert.Equal(t, env.Map{"FOO": "bar", "QUX": "x00"}, env.Environ())
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
		envs := envtest.Prepare(env.Map{"QUX": "x00"})
		defer envs.Restore()

		fsys := fstest.MapFS{"y": {Data: []byte("FOO=bar\nQUX=overload")}}
		assert.NoError(t, OverloadFS(fsys, "y"))
		assert.Equal(t, env.Map{"FOO": "bar", "QUX": "overload"}, env.Environ())
	})
}
