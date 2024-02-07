// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"os"
)

var _ Lookupper = (Map)(nil)

// Map represents a map of key value pairs.
type Map map[string]Value

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
func (m Map) Lookup(key string) (Value, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return "", errors.New(ErrNotFound)
}

// Load sets the environment variables from the Map using os.Setenv when they
// do not exist.
func (m Map) Load() error { return m.load(false) }

// Overload sets and overwrites the environment variables from the Map using
// os.Setenv
func (m Map) Overload() error { return m.load(true) }

func (m Map) load(overload bool) error {
	for k, v := range m {
		if err := set(k, v, overload); err != nil {
			return err
		}
	}
	return nil
}

func set(key string, val Value, overload bool) error {
	if !overload {
		_, exists := os.LookupEnv(key)
		if exists {
			return nil
		}
	}

	return Setenv(key, val)
}

// Clone returns a copy of the Map.
func (m Map) Clone() Map {
	clone := make(Map, len(m))
	clone.MergeValues(m)
	return clone
}
