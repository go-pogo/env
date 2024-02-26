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
	Environ() Map
}

type osEnv struct{}

// System returns an Environment which wraps the operating system's env related
// functions.
//
//	dec := NewDecoder(System())
func System() Environment { return new(osEnv) }

var environ = System()

// Use another Environment with this package. This is especially useful when
// testing.
//
//	myEnv := make(env.Map)
//	env.Use(myEnv)
//	env.Setenv("foo", "bar")
//
// To reset env to use the system's environment:
//
//	Use(System())
func Use(e Environment) { environ = e }

// Setenv sets the Value of the environment variable named by the key on the
// current Environment.
func Setenv(key string, val Value) error { return environ.Set(key, val) }

// Set sets the Value of the environment variable named by the key using
// os.Setenv.
func (osEnv) Set(key string, val Value) error {
	return errors.WithStack(os.Setenv(key, val.String()))
}

// Getenv retrieves the Value of the environment variable named by the key from
// the current Environment.
func Getenv(key string) Value { return environ.Get(key) }

// Get retrieves the Value of the environment variable named by the key using
// os.Getenv.
func (osEnv) Get(key string) Value { return Value(os.Getenv(key)) }

// LookupEnv retrieves the Value of the environment variable named by the key
// from the current Environment.
func LookupEnv(key string) (Value, bool) {
	v, err := environ.Lookup(key)
	return v, err != nil
}

// Lookup retrieves the Value of the environment variable named by the key using
// os.LookupEnv.
func (osEnv) Lookup(key string) (Value, error) {
	if v, ok := os.LookupEnv(key); ok {
		return Value(v), nil
	}
	return "", errors.New(ErrNotFound)
}

// Has indicates if the environment variable named by the key is present in the
// system's environment using os.LookupEnv.
func (osEnv) Has(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Environ returns a Map with the environment variables from the current
// Environment.
func Environ() Map { return environ.Environ() }

// Environ returns a Map with the environment variables of os.Environ.
func (osEnv) Environ() Map {
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
