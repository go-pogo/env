// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

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
			_ = LoadFS(nil, "", "")
		})
	})

	t.Run("load", func(t *testing.T) {
		envs := envtest.Prepare(env.Map{"FOO": "baz"})
		defer envs.Restore()

		fsys := fstest.MapFS{
			".env":      {Data: []byte("FOO=bar\nQUX=x00")},
			".env.test": {Data: []byte("QUX=XOO")},
			".env.dev":  {Data: []byte("QUX=")},
		}
		assert.NoError(t, LoadFS(fsys, "", Development))
		assert.Equal(t, env.Map{
			"FOO": "baz",
			"QUX": "",
		}, env.Environ())
	})
}

func TestOverloadFS(t *testing.T) {
	t.Run("nil fsys", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilFsys, func() {
			_ = OverloadFS(nil, "", "")
		})
	})

	t.Run("overload", func(t *testing.T) {
		envs := envtest.Prepare(env.Map{"FOO": "baz"})
		defer envs.Restore()

		fsys := fstest.MapFS{
			".env":      {Data: []byte("FOO=bar\nQUX=x00")},
			".env.test": {Data: []byte("QUX=XOO")},
			".env.dev":  {Data: []byte("QUX=devverdedevdev")},
		}
		assert.NoError(t, OverloadFS(fsys, "", Testing))
		assert.Equal(t, env.Map{
			"FOO": "bar",
			"QUX": "XOO",
		}, env.Environ())
	})
}
