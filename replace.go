// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/errors"
	"regexp"
)

const ErrCircularDependency errors.Msg = "circular dependency"

func ReplaceAll(m Map) (Map, error) {
	r := Replacer{
		src:    m,
		result: make(Map, len(m)),
		stack:  make([]string, 0, 2),
	}
	for k, v := range m {
		if err := r.handle(k, v); err != nil {
			return m, err
		}
	}
	return r.result, nil
}

type Replacer struct {
	src Lookupper
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
		src:    l,
		result: make(Map, n),
		stack:  make([]string, 0, 2),
	}
}

const chars = `[a-zA-Z0-9_-]+`

var matcher = regexp.MustCompile(`\$(` + chars + `|\{` + chars + `(:-.*)?\})`)

func (r *Replacer) handle(k string, v Value) error {
	if contains(r.stack, k) {
		return errors.New(ErrCircularDependency)
	}

	val := v.String()
	matches := matcher.FindAllStringSubmatchIndex(val, -1)
	if matches == nil {
		r.result[k] = v
		return nil
	}

	r.stack = append(r.stack, k)

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
				return err
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

	r.stack = r.stack[:len(r.stack)-1]
	r.result[k] = Value(val)

	return nil
}

func (r *Replacer) Lookup(k string) (Value, error) {
	if v, ok := r.result[k]; ok {
		return v, nil
	}

	v, err := r.src.Lookup(k)
	if err != nil {
		return v, err
	}

	err = r.handle(k, v)
	return r.result[k], err
}

func contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
