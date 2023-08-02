// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"io"
	"io/fs"
	"strings"
)

type Reader interface {
	Lookupper
	ReadAll() (Map, error)
}

// Deprecated: use Reader interface instead.
type LookupMapper = Reader

var _ Reader = new(reader)

type reader struct {
	scanner *Scanner
	found   Map
}

type FileReader struct {
	Reader
	file fs.File
}

func NewReader(r io.Reader) Reader {
	return &reader{
		scanner: NewScanner(r),
		found:   make(Map, 4),
	}
}

func NewStringReader(str string) Reader {
	return NewReader(strings.NewReader(str))
}

func NewFileReader(f fs.File) *FileReader {
	return &FileReader{
		Reader: NewReader(f),
		file:   f,
	}
}

// Open opens filename for reading using os.Open and returns a new *FileReader.
// It is the caller's responsibility to close the FileReader when finished.
// If there is an error, it will be of type *os.PathError.
func Open(filename string) (*FileReader, error) {
	return OpenFS(osFS{}, filename)
}

// OpenFS opens filename for reading from fsys and returns a new *FileReader.
// It is the caller's responsibility to close the FileReader when finished.
// If there is an error, it will be of type *os.PathError.
func OpenFS(fsys fs.FS, filename string) (*FileReader, error) {
	f, err := fsys.Open(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewFileReader(f), nil
}

// Close closes the underlying fs.File.
func (f *FileReader) Close() error { return f.file.Close() }

// Lookup continues reading and scanning the internal io.Reader until either
// EOF is reached or key is found. It will return the found value, ErrNotFound
// if not found, or an error if any has occurred while scanning.
func (r *reader) Lookup(key string) (Value, error) {
	if v, ok := r.found[key]; ok {
		return v, nil
	}

	v, found, err := r.scan(key)
	if !found && err == nil {
		err = errors.New(ErrNotFound)
	}
	return v, err
}

// ReadAll continues reading and scanning the internal io.Reader and returns a
// Map of all found environment variables when either EOF is reached or an
// error has occurred.
func (r *reader) ReadAll() (Map, error) {
	if _, _, err := r.scan(""); err != nil {
		return nil, err
	}

	return r.found, nil
}

// scan continues scanning the internal io.Reader until either EOF is reached or
// lookup is found. It will return the found value, a boolean indicating if the
// lookup was found and an error if any.
func (r *reader) scan(lookup string) (Value, bool, error) {
	for r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return "", false, err
		}

		env, err := r.scanner.NamedValue()
		if err != nil {
			return "", false, err
		}

		r.found[env.Name] = env.Value
		if lookup != "" && lookup == env.Name {
			return env.Value, true, nil
		}
	}

	return "", false, nil
}
