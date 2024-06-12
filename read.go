// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"io"

	"github.com/go-pogo/errors"
)

// Environment provides a Map of keys and values representing the environment.
type Environment interface {
	Environ() (Map, error)
}

type EnvironmentLookupper interface {
	Lookupper
	Environment
}

var _ EnvironmentLookupper = (*Reader)(nil)

type Reader struct {
	scanner *Scanner
	found   Map
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

// Environ continues reading and scanning the internal io.Reader and returns a
// Map of all found environment variables when either EOF is reached or an
// error has occurred.
func (r *Reader) Environ() (Map, error) {
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
