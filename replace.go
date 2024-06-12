// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"regexp"

	"github.com/go-pogo/errors"
)

const ErrCircularDependency errors.Msg = "circular dependency"

func ReplaceAll(m Map) (Map, error) {
	return replaceAll(m, &Replacer{
		lookupper: m,
		result:    make(Map, len(m)),
		stack:     make([]string, 0, 2),
	})
}

func replaceAll(m Map, r *Replacer) (Map, error) {
	for k, v := range m {
		if err := r.handle(k, v); err != nil {
			return m, err
		}
	}
	return r.result, nil
}

type Replacer struct {
	lookupper Lookupper
	// result contains already replaced values
	result Map
	// stack of keys that are being handled, used to detect circular dependencies
	stack []string
}

func NewReplacer(l Lookupper) *Replacer {
	n := 10
	if m, ok := l.(Map); ok {
		n = len(m)
	}

	return &Replacer{
		lookupper: l,
		result:    make(Map, n),
		stack:     make([]string, 0, 2),
	}
}

// Unwrap returns the original Lookupper that was wrapped by the Replacer.
func (r *Replacer) Unwrap() Lookupper { return r.lookupper }

func (r *Replacer) Lookup(k string) (Value, error) {
	if v, ok := r.result[k]; ok {
		return v, nil
	}

	v, err := r.lookupper.Lookup(k)
	if err != nil {
		return v, err
	}
	if err = r.handle(k, v); err != nil {
		return "", err
	}

	return r.result[k], nil
}

func (r *Replacer) handle(k string, v Value) error {
	if contains(r.stack, k) {
		return errors.New(ErrCircularDependency)
	}

	v, err := r.replace(k, v)
	if err != nil {
		return err
	}

	r.result[k] = v
	return nil
}

func (r *Replacer) Replace(v Value) (Value, error) { return r.replace("", v) }

const chars = `[a-zA-Z0-9_-]+`

var matcher = regexp.MustCompile(`\$(` + chars + `|\{` + chars + `(:-.*)?\})`)

func (r *Replacer) replace(k string, v Value) (Value, error) {
	val := v.String()
	matches := matcher.FindAllStringSubmatchIndex(val, -1)
	if matches == nil {
		return v, nil
	}

	if k != "" {
		r.stack = append(r.stack, k)
		defer func() {
			r.stack = r.stack[:len(r.stack)-1]
		}()
	}

	var offset int
	for _, m := range matches {
		var bash bool

		lookup := val[m[2]+offset : m[3]+offset]
		if lookup[0] == '{' {
			lookup = lookup[1 : len(lookup)-1]
			bash = true
		}

		var repl string
		if v, err := r.Lookup(lookup); err != nil {
			if !IsNotFound(err) {
				return v, err
			}
			if !bash || m[4] < 0 {
				// no bash style default value set
				continue
			}

			repl = val[m[4]+2+offset : m[5]+offset]
		} else {
			repl = v.String()
		}

		val = val[:m[0]+offset] + repl + val[m[1]+offset:]
		offset += len(repl) - (m[1] - m[0])
	}

	return Value(val), nil
}

func contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
