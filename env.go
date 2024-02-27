// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"os"
	"strings"
)

type osEnv struct{}

// System returns an EnvironmentLookupper which wraps the operating system's env related
// functions.
//
//	dec := NewDecoder(System())
func System() EnvironmentLookupper { return new(osEnv) }

// Setenv sets the Value of the environment variable named by the key using
// // os.Setenv.
func Setenv(key string, val Value) error {
	return errors.WithStack(os.Setenv(key, val.String()))
}

// Getenv retrieves the Value of the environment variable named by the key using
// os.Getenv.
func Getenv(key string) Value { return Value(os.Getenv(key)) }

// LookupEnv retrieves the Value of the environment variable named by the key
// using os.LookupEnv.
func LookupEnv(key string) (Value, bool) {
	v, ok := os.LookupEnv(key)
	return Value(v), ok
}

func (o osEnv) Lookup(key string) (Value, error) {
	v, ok := LookupEnv(key)
	if !ok {
		return "", errors.New(ErrNotFound)
	}
	return v, nil
}

// Environ returns a Map with the environment variables using os.Environ.
func Environ() Map {
	env := os.Environ()
	res := make(Map, len(env))

	for _, e := range env {
		if e[0] == '=' {
			continue
		}

		i := strings.IndexRune(e, '=')
		res[e[:i]] = Value(e[i+1:])
	}
	return res
}

func (o osEnv) Environ() (Map, error) { return Environ(), nil }
