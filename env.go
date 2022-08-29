// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package env handles reading and decoding of environment variables from files or
other data sources.
*/
package env

import (
	"os"
	"strings"
)

type Fallbacker interface {
	Fallback(f Lookupper)
}

type Lookupper interface {
	Lookup(key string) (Value, bool)
}

type LookupperFunc func(key string) (Value, bool)

func (f LookupperFunc) Lookup(key string) (Value, bool) { return f(key) }

// Map represents a map of key value pairs.
type Map map[string]Value

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

// Merge any map of strings into this Map. Existing keys in Map are overwritten
// with the value of the key in src.
func (m Map) Merge(src map[string]string) {
	for k, v := range src {
		m[k] = Value(v)
	}
}

// MergeValues merges a map of Value into this Map. Existing keys in Map m are
// overwritten with the value of the key in src.
func (m Map) MergeValues(src map[string]Value) {
	for k, v := range src {
		m[k] = v
	}
}

// Lookup retrieves the Value of the environment variable named by the key.
// If the key is present in Map, the value (which may be empty) is returned
// and the boolean is true. Otherwise, the returned value will be empty and the
// boolean is false.
func (m Map) Lookup(key string) (Value, bool) {
	v, ok := m[key]
	return v, ok
}

// Lookup retrieves the Value of the environment variable named by the key
// from any of the provided arguments.
// If the key is present the value (which may be empty) is returned and the
// boolean is true. Otherwise, the returned value will be empty and the boolean
// will be false.
func Lookup(key string, from ...Lookupper) (Value, bool) {
	for _, l := range from {
		if v, ok := l.Lookup(key); ok {
			return v, true
		}
	}
	return "", false
}

// LookupEnv retrieves the Value of the environment variable named by the key.
// It behaves similar to os.LookupEnv.
func LookupEnv(key string) (Value, bool) {
	v, ok := os.LookupEnv(key)
	return Value(v), ok
}

// Getenv retrieves the Value of the environment variable named by the key.
// It behaves similar to os.Getenv.
func Getenv(key string) Value { return Value(os.Getenv(key)) }
