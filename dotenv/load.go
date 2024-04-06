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
	return env.Load(Read(dir, ae))
}

// Overload sets and overwrites the environment variables from the active
// environment using env.Overload.
func Overload(dir string, ae ActiveEnvironment) error {
	return env.Overload(Read(dir, ae))
}

// LoadFS sets the environment variables from the active environment using
// env.Load.
func LoadFS(fsys fs.FS, dir string, ae ActiveEnvironment) error {
	return env.Load(ReadFS(fsys, dir, ae))
}

// OverloadFS sets and overwrites the environment variables from the active
// environment using env.Overload.
func OverloadFS(fsys fs.FS, dir string, ae ActiveEnvironment) error {
	return env.Overload(ReadFS(fsys, dir, ae))
}
