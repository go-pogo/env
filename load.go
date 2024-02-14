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

// LoadEnv sets the environment variable named by key if it does not exist.
func LoadEnv(key string, val Value) error {
	if _, exists := os.LookupEnv(key); !exists {
		return Setenv(key, val)
	}
	return nil
}

// OpenAndLoad reads environment variables from the filename and sets them using
// Setenv when they do not exist.
func OpenAndLoad(filename string) error { return OpenAndLoadFS(osFS{}, filename) }

// OpenAndOverload reads environment variables from the filename and overwrites
// them using Setenv when they already exist.
func OpenAndOverload(filename string) error { return OpenAndOverloadFS(osFS{}, filename) }

// OpenAndLoadFS reads environment variables from the filename in fsys and sets them
// using Setenv when they do not exist.
func OpenAndLoadFS(fsys fs.FS, filename string) error {
	return openAndLoad(fsys, filename, false)
}

// OpenAndOverloadFS reads environment variables from the filename in fsys and
// overwrites them using Setenv when they already exist.
func OpenAndOverloadFS(fsys fs.FS, filename string) error {
	return openAndLoad(fsys, filename, true)
}

// ReadAndLoad reads environment variables from r and sets them using Setenv
// when they do not exist.
func ReadAndLoad(r io.Reader) error { return readAndLoad(r, false) }

// ReadAndOverload reads environment variables from r and overwrites them using
// Setenv when they already exist.
func ReadAndOverload(r io.Reader) error { return readAndLoad(r, true) }

var _ fs.FS = (*osFS)(nil)

// osFS is a fs.FS compatible wrapper around os.Open.
type osFS struct{}

func (o osFS) Open(name string) (fs.File, error) { return os.Open(name) }

func openAndLoad(fsys fs.FS, filename string, overload bool) (err error) {
	f, err := fsys.Open(filename)
	if err != nil {
		return errors.WithStack(err)
	}

	defer errors.AppendFunc(&err, f.Close)
	return readAndLoad(f, overload)
}

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
