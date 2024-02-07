// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"io"
	"io/fs"
)

// ReadLookupper reads environment variables from any source.
type ReadLookupper interface {
	Lookupper
	ReadAll() (Map, error)
}

// ReadCloseLookupper reads environment variables from any source that needs to
// be closed when done.
type ReadCloseLookupper interface {
	ReadLookupper
	io.Closer
}

var _ ReadLookupper = (*Reader)(nil)

type Reader struct {
	scanner *Scanner
	found   Map
}

// reader prevents FileReader from needing to have a public *Reader
type reader = Reader

var _ ReadCloseLookupper = (*FileReader)(nil)

type FileReader struct {
	*reader
	file fs.File
}

// NewReader returns a Reader which looks up environment variables from
// the provided io.Reader r.
//
//	dec := NewDecoder(NewReader(r))
func NewReader(r io.Reader) *Reader {
	return &Reader{
		scanner: NewScanner(r),
		found:   make(Map, 4),
	}
}

// NewFileReader returns a Reader which looks up environment variables from
// the provided io.Reader r.
//
//	dec := NewDecoder(NewFileReader(f))
func NewFileReader(f fs.File) *FileReader {
	return &FileReader{
		reader: NewReader(f),
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
func (r *Reader) Lookup(key string) (Value, error) {
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
func (r *Reader) ReadAll() (Map, error) {
	if _, _, err := r.scan(""); err != nil {
		return nil, err
	}

	return r.found, nil
}

// scan continues scanning the internal io.Reader until either EOF is reached or
// lookup is found. It will return the found value, a boolean indicating if the
// lookup was found and an error if any.
func (r *Reader) scan(lookup string) (Value, bool, error) {
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
