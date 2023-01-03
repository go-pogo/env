// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"github.com/go-pogo/env"
	"github.com/go-pogo/errors"
	"io/fs"
	"os"
)

type Environment string

func (e Environment) String() string { return string(e) }

const (
	Development Environment = "dev"
	Testing     Environment = "test"
	Production  Environment = "prod"

	ErrNoFilesLoaded = "no files loaded"
)

func Read(dir string, ae Environment) *Reader {
	return NewReader(os.DirFS(dir), ae)
}

var _ env.LookupMapCloser = new(Reader)

type Reader struct {
	fsys  fs.FS
	files []*file
	found env.Map
}

type file struct {
	name      string
	reader    *env.FileReader
	notExists bool
}

func NewReader(fsys fs.FS, ae Environment) *Reader {
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

func (er *Reader) init(fsys fs.FS) {
	if er.files != nil {
		return
	}

	if fsys == nil {
		fsys = os.DirFS("")
	}
	er.fsys = fsys
	er.found = make(env.Map, 8)
	er.files = []*file{
		{name: ".env"},
		{name: ".env.local"},
	}
}

func (er *Reader) reader(f *file) (*env.FileReader, error) {
	if f.reader != nil || f.notExists {
		return f.reader, nil
	}

	fr, err := env.OpenFS(er.fsys, f.name)
	if err != nil {
		if os.IsNotExist(err) {
			f.notExists = true
			return nil, nil
		}
		return nil, err
	}

	//log.Printf("dotenv.Reader: reading from `%s`\n", f.name)
	f.reader = fr
	f.notExists = false
	return f.reader, nil
}

func (er *Reader) Lookup(key string) (env.Value, error) {
	er.init(nil)
	if v, ok := er.found[key]; ok {
		return v, nil
	}

	for i := len(er.files) - 1; i >= 0; i-- {
		r, err := er.reader(er.files[i])
		if err != nil {
			return "", err
		}
		if r == nil {
			continue
		}

		v, err := r.Lookup(key)
		if env.IsNotFound(err) {
			continue
		}
		return v, err

	}
	return "", errors.New(env.ErrNotFound)
}

func (er *Reader) Map() (env.Map, error) {
	er.init(nil)

	res := make(env.Map, 8)
	for _, f := range er.files {
		r, err := er.reader(f)
		if err != nil {
			return res, err
		}
		if r == nil {
			continue
		}

		m, err := r.Map()
		if err != nil {
			return nil, err
		}
		res.MergeValues(m)
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
