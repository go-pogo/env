// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"os"
	"strings"
)

type Environment interface {
	Lookupper
	Set(key string, val Value) error
	Get(key string) Value
	Has(key string) bool
}

// System returns an Environment which wraps the operating system's env related
// functions.
//
//	dec := NewDecoder(System())
func System() Environment { return new(osEnv) }

var environ Environment = new(osEnv)

// Setenv sets the Value of the environment variable named by the key using
// os.Setenv.
func Setenv(key string, val Value) error { return environ.Set(key, val) }

// Getenv retrieves the Value of the environment variable named by the key.
// It behaves similar to os.Getenv.
func Getenv(key string) Value { return environ.Get(key) }

// LookupEnv retrieves the Value of the environment variable named by the key
// using os.LookupEnv.
func LookupEnv(key string) (Value, bool) {
	v, err := environ.Lookup(key)
	return v, err != nil
}

// Environ returns a Map with the environment variables of os.Environ.
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

type osEnv struct{}

func (o *osEnv) Set(key string, val Value) error {
	return errors.WithStack(os.Setenv(key, val.String()))
}

func (o *osEnv) Get(key string) Value { return Value(os.Getenv(key)) }

func (o *osEnv) Has(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func (o *osEnv) Lookup(key string) (Value, error) {
	if v, ok := os.LookupEnv(key); ok {
		return Value(v), nil
	}
	return "", errors.New(ErrNotFound)
}
