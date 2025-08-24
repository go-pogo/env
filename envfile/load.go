// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"io/fs"

	"github.com/go-pogo/env"
	"github.com/go-pogo/env/internal/osfs"
)

// Load reads from filename and sets the environment variables using [env.Load].
func Load(filename string) error {
	return readAndLoad(osfs.FS{}, filename, env.Load)
}

// Overload reads from filename and sets and overwrites the environment
// variables using [env.Overload].
func Overload(filename string) error {
	return readAndLoad(osfs.FS{}, filename, env.Overload)
}

// LoadFS reads from filename and sets the environment variables using
// [env.Load].
func LoadFS(fsys fs.FS, filename string) error {
	if fsys == nil {
		panic(panicNilFsys)
	}

	return readAndLoad(fsys, filename, env.Load)
}

// OverloadFS reads from filename and sets and overwrites the environment
// variables using [env.Overload].
func OverloadFS(fsys fs.FS, filename string) error {
	if fsys == nil {
		panic(panicNilFsys)
	}

	return readAndLoad(fsys, filename, env.Overload)
}

func readAndLoad(fsys fs.FS, filename string, loadFn func(env.Mapper) error) error {
	r, err := OpenFS(fsys, filename)
	if err != nil {
		return err
	}
	return loadFn(r)
}
