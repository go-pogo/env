// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"github.com/go-pogo/env"
	"io/fs"
)

func Load(dir string, ae ActiveEnvironment) error {
	return readAndLoad(Read(dir, ae), false)
}

func Overload(dir string, ae ActiveEnvironment) error {
	return readAndLoad(Read(dir, ae), true)
}

func LoadFS(fsys fs.FS, ae ActiveEnvironment) error {
	return readAndLoad(ReadFS(fsys, ae), false)
}

func OverloadFS(fsys fs.FS, ae ActiveEnvironment) error {
	return readAndLoad(ReadFS(fsys, ae), true)
}

func readAndLoad(r *Reader, overload bool) error {
	defer r.Close()
	m, err := r.ReadAll()
	if err != nil {
		return err
	}
	if m, err = env.ReplaceAll(m); err != nil {
		return err
	}

	if overload {
		return m.Overload()
	} else {
		return m.Load()
	}
}
