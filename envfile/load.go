// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"github.com/go-pogo/env"
	"io/fs"
)

// Load reads from filename and sets the environment variables using env.Load.
func Load(filename string) error {
	return readAndLoad(nil, filename, env.Load)
}

// Overload reads from filename and sets and overwrites the environment
// variables using env.Overload.
func Overload(filename string) error {
	return readAndLoad(nil, filename, env.Overload)
}

// LoadFS reads from filename and sets the environment variables using env.Load.
func LoadFS(fsys fs.FS, filename string) error {
	return readAndLoad(fsys, filename, env.Load)
}

// OverloadFS reads from filename and sets and overwrites the environment
// variables using env.Overload.
func OverloadFS(fsys fs.FS, filename string) error {
	return readAndLoad(fsys, filename, env.Overload)
}

func readAndLoad(fsys fs.FS, filename string, loadFn func(env.Map) error) error {
	r, err := OpenFS(fsys, filename)
	if err != nil {
		return err
	}
	environ, err := r.Environ()
	if err != nil {
		return err
	}
	return loadFn(environ)
}