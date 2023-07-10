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

func Load(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return errors.WithStack(err)
	}

	defer f.Close()
	return load(f, false)
}

func LoadFS(fsys fs.FS, name string) error {
	f, err := fsys.Open(name)
	if err != nil {
		return errors.WithStack(err)
	}

	defer f.Close()
	return load(f, false)
}

func LoadFrom(r io.Reader) error { return load(r, false) }

func load(r io.Reader, overload bool) error {
	scan := NewScanner(r)
	for scan.Scan() {
		if err := scan.Err(); err != nil {
			return err
		}

		k, v, err := scan.KeyValue()
		if err != nil {
			return err
		}

		if err = set(k, v, overload); err != nil {
			return err
		}
	}
	return nil
}

func set(key string, val Value, overload bool) error {
	if !overload {
		_, exists := os.LookupEnv(key)
		if exists {
			return nil
		}
	}

	return Setenv(key, val)
}
