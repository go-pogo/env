// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import "strings"

// Load sets the system's environment variables with those from the Map when
// they do not exist.
func Load(envs Map) error {
	return load(envs, environ, false)
}

// Overload sets and overwrites the system's environment variables with those
// from the Map.
func Overload(envs Map) error {
	return load(envs, environ, true)
}

func load(envs Map, environ Environment, overload bool) error {
	if len(envs) == 0 {
		return nil
	}

	var r *Replacer
	if predictReplacerNeed(envs) {
		r = NewReplacer(Chain(envs, environ))
	}

	for k, v := range envs {
		if !overload && environ.Has(k) {
			continue
		}

		if r != nil {
			var err error
			v, err = r.Replace(v)
			if err != nil {
				return err
			}
		}

		if err := environ.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func predictReplacerNeed(m Map) bool {
	for _, v := range m {
		if strings.IndexRune(v.String(), '$') >= 0 {
			return true
		}
	}
	return false
}
