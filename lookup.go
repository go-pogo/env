// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import "github.com/go-pogo/errors"

type Lookupper interface {
	Lookup(key string) (val Value, err error)
}

type LookupperFunc func(key string) (Value, error)

func (f LookupperFunc) Lookup(key string) (Value, error) { return f(key) }

type LookupMapper interface {
	Lookupper
	Map() (Map, error)
}

// ErrNotFound is returned when a Lookup call cannot find a matching value.
const ErrNotFound errors.Msg = "not found"

// IsNotFound tests whether the provided error is ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// Lookup retrieves the Value of the environment variable named by the key
// from any of the provided Lookupper(s).
// If the key is present the value (which may be empty) is returned and the
// error is nil. Otherwise, the returned value will be empty and the error
// ErrNotFound.
func Lookup(key string, from ...Lookupper) (Value, error) {
	for _, l := range from {
		if v, err := l.Lookup(key); IsNotFound(err) {
			continue
		} else {
			return v, err
		}
	}
	return "", errors.New(ErrNotFound)
}

var _ Lookupper = new(chain)

type chain []Lookupper

// Chain multiple Lookupper(s) to lookup keys from.
func Chain(l ...Lookupper) Lookupper {
	if n := len(l); n == 1 {
		return l[0]
	}
	if c, ok := l[0].(chain); ok {
		c = append(c, l[1:]...)
		return c
	}
	return chain(l)
}

func (c chain) Lookup(key string) (Value, error) { return Lookup(key, c...) }