// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"github.com/go-pogo/env"
	"github.com/go-pogo/errors"
	"io/fs"
	"os"
	"path"
)

const ErrNoFilesLoaded = "no files loaded"

var (
	_ env.Lookupper = new(Reader)
	_ env.AllReader = new(Reader)
)

type Reader struct {
	fsys  fs.FS
	files []*file
	found env.Map

	// ReplaceVars
	ReplaceVars bool
}

type file struct {
	name      string
	reader    *env.FileReader
	notExists bool
}

func Read(dir string, ae Environment) *Reader {
	var fsys fs.FS
	if dir != "" {
		fsys = os.DirFS(dir)
	}

	return ReadFS(fsys, ae)
}

func ReadFS(fsys fs.FS, ae Environment) *Reader {
	var er Reader
	er.init(fsys)

	if ae != "" {
		er.files = append(er.files,
			&file{name: ".env." + ae.String()},
			&file{name: ".env." + ae.String() + ".local"},
		)
	}

	return &er
}

// Deprecated: use ReadFS instead.
func NewReader(fsys fs.FS, ae Environment) *Reader { return ReadFS(fsys, ae) }

func (er *Reader) init(fsys fs.FS) {
	if er.files != nil {
		return
	}
	if fsys == nil {
		dir, err := os.Getwd()
		if err != nil {
			if dir, err = os.Executable(); err == nil {
				dir = path.Dir(dir)
			}
		}
		fsys = os.DirFS(dir)
	}

	er.fsys = fsys
	er.found = make(env.Map, 8)
	er.files = []*file{
		{name: ".env"},
		{name: ".env.local"},
	}

	er.ReplaceVars = true
}

func (er *Reader) reader(f *file) (*env.FileReader, error) {
	if f.reader != nil || f.notExists {
		return f.reader, nil
	}

	fr, err := env.OpenFS(er.fsys, f.name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			f.notExists = true
			return nil, nil
		}
		return nil, err
	}

	f.reader = fr
	f.notExists = false
	return f.reader, nil
}

func (er *Reader) Lookup(key string) (env.Value, error) {
	er.init(nil)
	if v, ok := er.found[key]; ok {
		return v, nil
	}

	var anyLoaded bool
	for i := len(er.files) - 1; i >= 0; i-- {
		r, err := er.reader(er.files[i])
		if err != nil {
			return "", err
		}
		if r == nil {
			anyLoaded = true
			continue
		}

		v, err := r.Lookup(key)
		if env.IsNotFound(err) {
			continue
		}
		return v, err

	}
	if !anyLoaded {
		return "", errors.New(ErrNoFilesLoaded)
	}

	return "", errors.New(env.ErrNotFound)
}

func (er *Reader) ReadAll() (env.Map, error) {
	er.init(nil)
	var anyLoaded bool

	res := make(env.Map, 8)
	for _, f := range er.files {
		r, err := er.reader(f)
		if err != nil {
			return res, err
		}

		m, err := r.ReadAll()
		if err != nil {
			return nil, err
		}

		res.MergeValues(m)
		anyLoaded = true
	}
	if !anyLoaded {
		return nil, errors.New(ErrNoFilesLoaded)
	}

	er.found.MergeValues(res)
	return res, nil
}

func (er *Reader) Close() error {
	var err error
	for _, f := range er.files {
		if f.reader != nil {
			f.notExists = false
			errors.Append(&err, f.reader.Close())
		}
	}
	return err
}
