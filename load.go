// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"strings"
)

// Load sets the system's environment variables with those from the Map when
// they do not exist.
func Load(envs Environment) error {
	return load(envs, false)
}

// Overload sets and overwrites the system's environment variables with those
// from the Map.
func Overload(envs Environment) error {
	return load(envs, true)
}

func load(envs Environment, overload bool) (err error) {
	var m Map
	if em, ok := envs.(Map); ok {
		m = em
	} else if m, err = envs.Environ(); err != nil {
		return errors.WithStack(err)
	}

	if len(m) == 0 {
		return nil
	}

	var r *Replacer
	if predictReplacerNeed(m) {
		r = NewReplacer(Chain(m, System()))
	}

	for k, v := range m {
		if _, has := LookupEnv(k); has && !overload {
			continue
		}

		if r != nil {
			v, err = r.Replace(v)
			if err != nil {
				return err
			}
		}

		if err = Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func predictReplacerNeed(m Map) bool {
	for _, v := range m {
		if strings.ContainsRune(v.String(), '$') {
			return true
		}
	}
	return false
}
