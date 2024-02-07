// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
)

// ErrNotFound is returned when a Lookup call cannot find a matching key.
const ErrNotFound errors.Msg = "not found"

// IsNotFound tests whether the provided error is ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(errors.Unembed(err), ErrNotFound)
}

type Lookupper interface {
	// Lookup retrieves the Value of the environment variable named by the key.
	// It must return an ErrNotFound error if the key is not present.
	Lookup(key string) (val Value, err error)
}

type LookupperFunc func(key string) (Value, error)

func (f LookupperFunc) Lookup(key string) (Value, error) { return f(key) }

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

var _ Lookupper = (chainLookupper)(nil)

type chainLookupper []Lookupper

// Chain multiple Lookupper(s) to lookup keys from.
func Chain(l ...Lookupper) Lookupper {
	res, chained := chain(l...)
	if !chained && res == nil {
		return make(chainLookupper, 0)
	}
	return res
}

func chain(lookuppers ...Lookupper) (Lookupper, bool) {
	if n := len(lookuppers); n == 1 {
		return lookuppers[0], false
	}

	var res chainLookupper
	if c, ok := lookuppers[0].(chainLookupper); ok {
		res = c
		lookuppers = lookuppers[1:]
	} else {
		res = make(chainLookupper, 0, len(lookuppers))
	}
	for _, l := range lookuppers {
		if l != nil {
			res = append(res, l)
		}
	}
	return res, true
}

func (c chainLookupper) Lookup(key string) (Value, error) { return Lookup(key, c...) }
