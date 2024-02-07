// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"io"
	"io/fs"
	"os"
)

var _ fs.FS = (*osFS)(nil)

// osFS is a fs.FS compatible wrapper around os.Open.
type osFS struct{}

func (o osFS) Open(name string) (fs.File, error) { return os.Open(name) }

// Load reads environment variables from the filename and sets them using
// Setenv when they do not exist.
func Load(filename string) error { return LoadFS(osFS{}, filename) }

// Overload reads environment variables from the filename and overwrites them
// using Setenv when they already exist.
func Overload(filename string) error { return OverloadFS(osFS{}, filename) }

// LoadFS reads environment variables from the filename in fsys and sets them
// using Setenv when they do not exist.
func LoadFS(fsys fs.FS, filename string) error {
	return openAndLoad(fsys, filename, false)
}

// OverloadFS reads environment variables from the filename in fsys and
// overwrites them using Setenv when they already exist.
func OverloadFS(fsys fs.FS, filename string) error {
	return openAndLoad(fsys, filename, true)
}

func openAndLoad(fsys fs.FS, filename string, overload bool) (err error) {
	f, err := fsys.Open(filename)
	if err != nil {
		return errors.WithStack(err)
	}

	defer errors.AppendFunc(&err, f.Close)
	return readAndLoad(f, overload)
}

// LoadFrom reads environment variables from r and sets them using Setenv when
// they do not exist.
func LoadFrom(r io.Reader) error { return readAndLoad(r, false) }

// OverloadFrom reads environment variables from r and overwrites them using
// Setenv when they already exist.
func OverloadFrom(r io.Reader) error { return readAndLoad(r, true) }

func readAndLoad(r io.Reader, overload bool) error {
	m, err := NewReader(r).ReadAll()
	if err != nil {
		return err
	}
	if m, err = ReplaceAll(m); err != nil {
		return err
	}

	return m.load(overload)
}
