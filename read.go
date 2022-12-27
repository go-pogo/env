// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"io"
	"io/fs"
	"os"
)

type LookupMapCloser interface {
	LookupMapper
	io.Closer
}

var (
	_ LookupMapper    = new(Reader)
	_ LookupMapCloser = new(FileReader)
)

type Reader struct {
	scanner Scanner
	found   Map
}

type FileReader struct {
	reader *Reader
	file   fs.File
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		scanner: NewScanner(r),
		found:   make(Map, 4),
	}
}

func NewFileReader(f fs.File) *FileReader {
	return &FileReader{
		reader: NewReader(f),
		file:   f,
	}
}

func Open(name string) (*FileReader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewFileReader(f), nil
}

func (f *FileReader) Close() error { return f.file.Close() }

func (f *FileReader) Lookup(key string) (Value, error) {
	return f.reader.Lookup(key)
}

// Lookup a value by scanning the internal io.Reader.
func (r *Reader) Lookup(key string) (Value, error) {
	if v, ok := r.found[key]; ok {
		return v, nil
	}

	v, ok, err := r.scan(key)
	if !ok && err == nil {
		err = errors.New(ErrNotFound)
	}
	return v, err
}

func (f *FileReader) Map() (Map, error) { return f.reader.Map() }

// Map returns a Map of all found environment variables.
func (r *Reader) Map() (Map, error) {
	if _, _, err := r.scan(""); err != nil {
		return nil, err
	}

	return r.found, nil
}

func (r *Reader) scan(lookup string) (Value, bool, error) {
	for r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return "", false, err
		}

		k, val, err := parse(r.scanner.Text())
		if err != nil {
			return "", false, err
		}

		r.found[k] = val
		if lookup != "" && lookup == k {
			return val, true, nil
		}
	}

	return "", false, nil
}
