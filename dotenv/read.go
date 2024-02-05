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

const ErrNoFilesLoaded errors.Msg = "no files loaded"

var (
	_ env.Lookupper = new(Reader)
	_ env.ReadAller = new(Reader)
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

// Read reads .env files from dir depending on the provided ActiveEnvironment.
//
//	var cfg MyConfig
//	env.NewDecoder(dotenv.ReadAll("./", dotenv.Development)).Decode(&cfg)
func Read(dir string, ae ActiveEnvironment) *Reader {
	var fsys fs.FS
	if dir != "" {
		fsys = os.DirFS(dir)
	}

	return ReadFS(fsys, ae)
}

// ReadFS reads .env files from fsys.
func ReadFS(fsys fs.FS, ae ActiveEnvironment) *Reader {
	var r Reader
	r.init(fsys)

	if ae != "" {
		r.files = append(r.files,
			&file{name: ".env." + ae.String()},
			&file{name: ".env." + ae.String() + ".local"},
		)
	}

	return &r
}

func (r *Reader) init(fsys fs.FS) {
	if r.files != nil {
		return
	}
	if fsys == nil {
		fsys = os.DirFS(getwd())
	}

	r.fsys = fsys
	r.found = make(env.Map, 8)
	r.files = []*file{
		{name: ".env"},
		{name: ".env.local"},
	}

	r.ReplaceVars = true
}

func (r *Reader) reader(f *file) (*env.FileReader, error) {
	if f.reader != nil || f.notExists {
		return f.reader, nil
	}

	fr, err := env.OpenFS(r.fsys, f.name)
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

func (r *Reader) Lookup(key string) (env.Value, error) {
	r.init(nil)
	if v, ok := r.found[key]; ok {
		return v, nil
	}

	var anyLoaded bool
	for i := len(r.files) - 1; i >= 0; i-- {
		fr, err := r.reader(r.files[i])
		if err != nil {
			return "", err
		}
		if fr == nil {
			anyLoaded = true
			continue
		}

		v, err := fr.Lookup(key)
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

func (r *Reader) ReadAll() (env.Map, error) {
	r.init(nil)
	var anyLoaded bool

	res := make(env.Map, 8)
	for _, f := range r.files {
		fr, err := r.reader(f)
		if err != nil {
			return res, err
		}
		if fr == nil {
			continue
		}

		m, err := fr.ReadAll()
		if err != nil {
			return nil, err
		}

		res.MergeValues(m)
		anyLoaded = true
	}
	if !anyLoaded {
		return nil, errors.New(ErrNoFilesLoaded)
	}

	r.found.MergeValues(res)
	return res, nil
}

func (r *Reader) Close() error {
	var err error
	for _, f := range r.files {
		if f.reader != nil {
			f.notExists = false
			errors.AppendFunc(&err, f.reader.Close)
		}
	}
	return err
}

// todo: test getwd vs "./"
func getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		if dir, err = os.Executable(); err == nil {
			dir = path.Dir(dir)
		} else {
			dir = "./"
		}
	}
	return dir
}
