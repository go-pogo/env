// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"github.com/go-pogo/env"
	"io/fs"
)

// Load sets the environment variables from the active environment using
// env.Load.
func Load(dir string, ae ActiveEnvironment) error {
	return readAndLoad(ReadFS(nil, dir, ae), env.Load)
}

// Overload sets and overwrites the environment variables from the active
// environment using env.Overload.
func Overload(dir string, ae ActiveEnvironment) error {
	return readAndLoad(ReadFS(nil, dir, ae), env.Overload)
}

// LoadFS sets the environment variables from the active environment using
// env.Load.
func LoadFS(fsys fs.FS, dir string, ae ActiveEnvironment) error {
	return readAndLoad(ReadFS(fsys, dir, ae), env.Load)
}

// OverloadFS sets and overwrites the environment variables from the active
// environment using env.Overload.
func OverloadFS(fsys fs.FS, dir string, ae ActiveEnvironment) error {
	return readAndLoad(ReadFS(fsys, dir, ae), env.Overload)
}

func readAndLoad(r *Reader, loadFn func(env.Map) error) error {
	environ, err := r.Environ()
	if err != nil {
		return err
	}
	return loadFn(environ)
}
