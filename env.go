// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package env handles reading and decoding of environment variables from files or
other data sources.
*/
package env

import (
	"github.com/go-pogo/errors"
	"os"
	"strings"
)

// Getenv retrieves the Value of the environment variable named by the key.
// It behaves similar to os.Getenv.
func Getenv(key string) Value { return Value(os.Getenv(key)) }

// LookupEnv retrieves the Value of the environment variable named by the key.
// It behaves similar to os.LookupEnv.
func LookupEnv(key string) (Value, bool) {
	v, ok := os.LookupEnv(key)
	return Value(v), ok
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

func EnvironLookup() Lookupper {
	return LookupperFunc(func(key string) (Value, error) {
		if v, ok := LookupEnv(key); ok {
			return v, nil
		}
		return "", errors.New(ErrNotFound)
	})
}
