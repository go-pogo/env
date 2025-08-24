// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
)

// Mapper provides a [Map] of keys and values representing the environment.
type Mapper interface {
	Environ() (Map, error)
}

var _ LookupMapper = (Map)(nil)

// Map represents a map of key value pairs.
type Map map[string]Value

// Lookup retrieves the [Value] of the environment variable named by the key.
// If the key is present in [Map], the value (which may be empty) is returned
// and the boolean is true. Otherwise, the returned value will be empty and the
// boolean is false.
func (m Map) Lookup(key string) (Value, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return "", errors.New(ErrNotFound)
}

// Merge any map of strings into this [Map]. Existing keys in the [Map] are
// overwritten with the value of the key in src.
func (m Map) Merge(src map[string]string) {
	for k, v := range src {
		m[k] = Value(v)
	}
}

// MergeValues merges a map of [Value] into this [Map]. Existing keys in the
// [Map] are overwritten with the value of the key in src.
func (m Map) MergeValues(src map[string]Value) {
	for k, v := range src {
		m[k] = v
	}
}

// Environ returns a copy of the [Map].
func (m Map) Environ() (Map, error) {
	clone := make(Map, len(m))
	clone.MergeValues(m)
	return clone, nil
}
