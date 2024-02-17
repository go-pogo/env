// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"github.com/go-pogo/env"
	"github.com/go-pogo/errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const ErrNoFilesLoaded errors.Msg = "no files loaded"

var (
	_ env.EnvironLookupper = (*Reader)(nil)
	_ io.Closer            = (*Reader)(nil)
)

type Reader struct {
	reader
	replacer *env.Replacer

	// ReplaceVars
	ReplaceVars bool
}

// Read reads .env files from dir depending on the provided ActiveEnvironment.
//
//	var cfg MyConfig
//	dec := env.NewDecoder(dotenv.Read("./", dotenv.Development))
//	dec.Decode(&cfg)
func Read(dir string, ae ActiveEnvironment) *Reader {
	return ReadFS(nil, dir, ae)
}

// ReadFS reads .env files at dir from fsys.
func ReadFS(fsys fs.FS, dir string, ae ActiveEnvironment) *Reader {
	r := Reader{ReplaceVars: true}
	r.init(fsys, dir)
	//goland:noinspection GoUnhandledErrorResult
	defer r.Close()

	if ae != "" {
		r.files = append(r.files,
			&file{name: ".env." + ae.String()},
			&file{name: ".env." + ae.String() + ".local"},
		)
	}
	return &r
}

var _ env.Lookupper = (*reader)(nil)

type reader struct {
	fsys  fs.FS
	dir   string
	files []*file
	found env.Map
}

func (r *reader) init(fsys fs.FS, dir string) {
	if r.files != nil {
		return
	}
	if fsys == nil {
		fsys = osFS{}
	}

	r.found = make(env.Map, 8)
	r.fsys = fsys
	r.dir = dir
	r.files = []*file{
		{name: ".env"},
		{name: ".env.local"},
	}
}

type file struct {
	name      string
	reader    *env.FileReader
	notExists bool
}

func (r *reader) fileReader(f *file) (*env.FileReader, bool, error) {
	if f.reader != nil || f.notExists {
		return f.reader, !f.notExists, nil
	}

	filename := f.name
	if r.dir != "" {
		filename = filepath.Join(r.dir, filename)
	}

	fr, err := env.OpenFS(r.fsys, filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			f.notExists = true
			return nil, !f.notExists, nil
		}
		return nil, !f.notExists, err
	}

	f.reader = fr
	f.notExists = false
	return f.reader, !f.notExists, nil
}

func (r *Reader) Lookup(key string) (env.Value, error) {
	if !r.ReplaceVars {
		if r.replacer != nil {
			r.replacer = nil
		}
		return r.reader.Lookup(key)
	}

	if r.replacer == nil {
		r.replacer = env.NewReplacer(&r.reader)
	}
	return r.replacer.Lookup(key)
}

func (r *reader) Lookup(key string) (env.Value, error) {
	r.init(nil, "")
	if v, ok := r.found[key]; ok {
		return v, nil
	}

	var anyLoaded bool
	for i := len(r.files) - 1; i >= 0; i-- {
		fr, exists, err := r.fileReader(r.files[i])
		anyLoaded = anyLoaded || exists
		if err != nil {
			return "", err
		}
		if fr == nil {
			continue
		}

		v, err := fr.Lookup(key)
		if err != nil {
			if env.IsNotFound(err) {
				continue
			}
			return v, err
		}
		return v, nil
	}
	if !anyLoaded {
		return "", errors.New(ErrNoFilesLoaded)
	}

	return "", errors.New(env.ErrNotFound)
}

// Environ reads and returns all environment variables from the loaded .env
// files.
func (r *reader) Environ() (env.Map, error) {
	r.init(nil, "")
	var anyLoaded bool

	res := make(env.Map, 8)
	for _, f := range r.files {
		fr, exists, err := r.fileReader(f)
		anyLoaded = anyLoaded || exists
		if err != nil {
			return res, err
		}
		if fr == nil {
			continue
		}

		m, err := fr.Environ()
		if err != nil {
			return nil, err
		}

		res.MergeValues(m)
	}
	if !anyLoaded {
		return nil, errors.New(ErrNoFilesLoaded)
	}

	r.found.MergeValues(res)
	return res, nil
}

func (r *reader) Close() error {
	var err error
	for _, f := range r.files {
		if f.reader != nil {
			f.notExists = false
			errors.AppendFunc(&err, f.reader.Close)
		}
	}
	return err
}

var _ fs.FS = (*osFS)(nil)

// osFS is a fs.FS compatible wrapper around os.Open.
type osFS struct{}

func (o osFS) Open(name string) (fs.File, error) { return os.Open(name) }
