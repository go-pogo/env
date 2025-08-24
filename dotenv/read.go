// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/go-pogo/env"
	"github.com/go-pogo/env/envfile"
	"github.com/go-pogo/env/internal/osfs"
	"github.com/go-pogo/errors"
)

type NoFilesLoadedError struct {
	FS  fs.FS
	Dir string
}

func (e *NoFilesLoadedError) Error() string {
	return fmt.Sprintf("no .env files loaded from directory `%s`", e.Dir)
}

var (
	_ env.LookupMapper = (*Reader)(nil)
	_ io.Closer        = (*Reader)(nil)
)

// A Reader reads .env files from a filesystem and provides the mechanism to
// lookup environment variables. Its zero value is ready to use and reads from
// the current working directory.
type Reader struct {
	fsys  fsJoiner
	dir   string
	files []*file
	found env.Map
}

// Read reads .env files from dir, depending on the provided ActiveEnvironment.
//
//	var cfg MyConfig
//	dec := env.NewDecoder(dotenv.Read("./", dotenv.Development))
//	dec.Decode(&cfg)
func Read(dir string, ae ActiveEnvironment) *Reader {
	return newReader(nil, dir, ae)
}

const panicNilFsys = "dotenv: fs.FS must not be nil"

// ReadFS reads .env files at dir from fsys, depending on the provided
// ActiveEnvironment.
func ReadFS(fsys fs.FS, dir string, ae ActiveEnvironment) *Reader {
	if fsys == nil {
		panic(panicNilFsys)
	}
	return newReader(fsys, dir, ae)
}

func newReader(fsys fs.FS, dir string, ae ActiveEnvironment) *Reader {
	var r Reader
	r.init(fsys, dir)

	//nolint:errcheck
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

func (r *Reader) init(fsys fs.FS, dir string) {
	if r.files != nil {
		return
	}
	if fsys == nil {
		r.fsys = osfs.FS{}
	} else {
		r.fsys = joinerFS{fsys}
	}

	r.found = make(env.Map, 8)
	r.dir = dir
	r.files = []*file{
		{name: ".env"},
		{name: ".env.local"},
	}
}

type file struct {
	name      string
	reader    *envfile.Reader
	notExists bool
}

func (r *Reader) fileReader(f *file) (*envfile.Reader, bool, error) {
	if f.reader != nil || f.notExists {
		return f.reader, !f.notExists, nil
	}

	filename := f.name
	if r.dir != "" {
		filename = r.fsys.JoinFilePath(r.dir, filename)
	}

	fr, err := envfile.OpenFS(r.fsys, filename)
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

// Lookup key by reading from .env files.
func (r *Reader) Lookup(key string) (env.Value, error) {
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
		return "", errors.WithStack(&NoFilesLoadedError{FS: r.fsys, Dir: r.dir})
	}

	return "", errors.New(env.ErrNotFound)
}

// Environ reads and returns all environment variables from the loaded .env
// files.
func (r *Reader) Environ() (env.Map, error) {
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
		return nil, errors.WithStack(&NoFilesLoadedError{FS: r.fsys, Dir: r.dir})
	}

	r.found.MergeValues(res)
	return res, nil
}

// Close closes all opened .env files.
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

type fsJoiner interface {
	fs.FS
	JoinFilePath(elem ...string) string
}

var _ fsJoiner = (*joinerFS)(nil)

type joinerFS struct{ fs.FS }

func (joinerFS) JoinFilePath(elem ...string) string { return path.Join(elem...) }
