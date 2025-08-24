// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"io"
	"io/fs"

	"github.com/go-pogo/env"
	"github.com/go-pogo/env/internal/osfs"
	"github.com/go-pogo/errors"
)

var (
	_ env.LookupMapper = (*Reader)(nil)
	_ io.Closer        = (*Reader)(nil)
)

// reader prevents Reader from needing to have a public *Reader
type reader = env.Reader

type Reader struct {
	*reader
	file fs.File
}

// NewReader returns a [Reader] which looks up environment variables from
// the provided [fs.File].
//
//	dec := env.NewDecoder(envfile.NewReader(file))
func NewReader(f fs.File) *Reader {
	if f == nil {
		panic(panicNilFile)
	}
	return &Reader{
		reader: env.NewReader(f),
		file:   f,
	}
}

// Open opens filename for reading using [os.Open] and returns a new [Reader].
// It is the caller's responsibility to close the [Reader] when finished.
// If there is an error, it will be of type *[os.PathError].
func Open(filename string) (*Reader, error) {
	return OpenFS(osfs.FS{}, filename)
}

// OpenFS opens filename for reading from fsys and returns a new [Reader].
// It is the caller's responsibility to close the [Reader] when finished.
// If there is an error, it will be of type *[os.PathError].
func OpenFS(fsys fs.FS, filename string) (*Reader, error) {
	if fsys == nil {
		panic(panicNilFsys)
	}

	f, err := fsys.Open(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewReader(f), nil
}

// Close closes the underlying [fs.File].
func (f *Reader) Close() error {
	return errors.WithStack(f.file.Close())
}
