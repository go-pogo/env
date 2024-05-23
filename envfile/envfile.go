// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package envfile provides tools to read and load environment variables from
// files.
package envfile

import (
	"github.com/go-pogo/errors"
	"path"
	"path/filepath"
	"runtime"
)

const (
	panicNilFile = "envfile: file must not be nil"
	panicNilFsys = "envfile: fs.FS must not be nil"
)

// Generate encodes and writes an env file in dir based on the provided src.
// It is meant to be used with go generate to create .env files based on the
// project's config(s).
func Generate(dir, filename string, src any) error {
	if dir == "" {
		_, dir, _, _ = runtime.Caller(1)
		dir = filepath.Dir(dir)
	}

	if !path.IsAbs(filename) {
		filename = filepath.Join(dir, ".env")
	}

	enc, err := Create(filename)
	if err != nil {
		return errors.WithStack(err)
	}

	enc.TakeValues = true
	defer errors.AppendFunc(&err, enc.Close)

	if err = enc.Encode(src); err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}
