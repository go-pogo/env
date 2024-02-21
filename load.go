// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

// LoadEnv sets the environment variable named by key if it does not exist.
// Note: An OverloadEnv function does not exist, because its functionality would
// be the same as Setenv.
func LoadEnv(key string, val Value) error {
	if environ.Has(key) {
		return nil
	}
	return environ.Set(key, val)
}

// Load sets the environment variables from the Map using Environment.Set when
// they do not exist.
func Load(environ Map) error {
	if len(environ) == 0 {
		return nil
	}

	return (&Loader{ReplaceVars: true}).Load(environ)
}

// Overload sets and overwrites the environment variables from the Map using
// Environment.Set.
func Overload(environ Map) error {
	if len(environ) == 0 {
		return nil
	}

	return (&Loader{Overload: true, ReplaceVars: true}).Load(environ)
}

type Loader struct {
	environ Environment

	Overload    bool
	ReplaceVars bool
}

func NewLoader(dest Environment) *Loader {
	return &Loader{
		environ:     dest,
		ReplaceVars: true,
	}
}

// Load sets the environment variables from the Map using os.Setenv when they
// do not exist.
func (l *Loader) Load(m Map) error {
	if len(m) == 0 {
		return nil
	}

	if l.environ == nil {
		l.environ = environ
	}

	var lookupper Lookupper = m
	if l.ReplaceVars {
		lookupper = NewReplacer(Chain(m, l.environ))
	}

	for k := range m {
		if !l.Overload && l.environ.Has(k) {
			continue
		}

		v, err := lookupper.Lookup(k)
		if err != nil {
			if !IsNotFound(err) {
				return err
			}
			continue
		}
		if err = l.environ.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}
