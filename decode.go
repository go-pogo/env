// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/parseval"
)

const ErrStructPointerExpected errors.Msg = "expected a pointer to a struct"

var parser *parseval.Parser

func init() {
	parser = parseval.NewParser()
	parser.Register(
		reflect.TypeOf((*Unmarshaler)(nil)).Elem(),
		func(val parseval.Value, dest interface{}) error {
			return dest.(Unmarshaler).UnmarshalEnv(val.Bytes())
		},
	)
}

// Unmarshaler is the interface implemented by types that can unmarshal a
// textual representation of themselves.
// It is similar to encoding.TextUnmarshaler.
type Unmarshaler interface {
	UnmarshalEnv([]byte) error
}

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(NewReader(bytes.NewBuffer(data))).Decode(v)
}

type Decoder struct {
	src Lookupper
	err error

	// TagsOnly ignores fields that do not have an `env` tag when set to true.
	TagsOnly bool

	// ReplaceVars
	ReplaceVars bool
}

// NewDecoder returns a new Decoder that scans io.Reader r for environment
// variables and parses them.
func NewDecoder(l Lookupper) *Decoder {
	return &Decoder{
		src: l,

		ReplaceVars: true,
	}
}

func (d *Decoder) Decode(v interface{}) error {
	if v == nil || reflect.TypeOf(v).Kind() != reflect.Ptr {
		return ErrStructPointerExpected
	}

	val := reflect.ValueOf(v)
	if underlyingKind(val) != reflect.Struct {
		return ErrStructPointerExpected
	}

	return d.decodeStruct(val, nil)
}

const panicPtr = "parseval.Indirect should always resolve ptr values; this is a bug!"

func (d *Decoder) decodeStruct(pv reflect.Value, p path) error {
	pv = indirect(pv)
	pt := pv.Type()

	for i := 0; i < pv.NumField(); i++ {
		field, rv := pt.Field(i), pv.Field(i)

		switch underlyingKind(rv) {
		case reflect.Invalid, reflect.Uintptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:
			continue

		case reflect.Pointer:
			panic(panicPtr)

		case reflect.Struct:
			if !parser.HasFunc(rv.Type()) {
				// continue traversing the struct...
				if err := d.decodeStruct(rv, p.extend(field.Name)); err != nil {
					return err
				} else {
					continue
				}
			}

			// parser supports the struct as a type, let it handle it further
			fallthrough

		default:
			if err := d.decodeField(field, rv, p.extend(field.Name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Decoder) decodeField(field reflect.StructField, rv reflect.Value, path path) error {
	t, found := field.Tag.Lookup(envtag.Key)
	if !found && d.TagsOnly {
		return nil
	}

	tag := envtag.ParseTag(t)
	if tag.Ignore {
		return nil
	}

	if tag.Name == "" {
		if tag.NoPrefix {
			tag.Name = path.last()
		} else {
			tag.Name = path.join()
		}
	}

	val, err := d.src.Lookup(tag.Name)
	if IsNotFound(err) {
		if def := field.Tag.Get("default"); def == "" {
			return nil
		} else {
			val = Value(def)
		}
	} else if err != nil {
		return err
	}

	return parser.Parse(val, rv)
}

func underlyingKind(rv reflect.Value) reflect.Kind {
	k := rv.Kind()
	for k == reflect.Ptr {
		rv = rv.Elem()
		k = rv.Kind()
	}
	return k
}

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// create a pointer to the type v points to
			ptr := reflect.New(v.Type().Elem())
			v.Set(ptr)
		}

		v = v.Elem()
	}
	return v
}

type path []string

func (p path) last() string { return p[len(p)-1] }

func (p path) join() string { return strings.Join(p, "_") }

func (p path) extend(s string) path {
	var x path
	if n := len(p); n == 0 {
		x = make(path, 0, 2)
	} else {
		x = make(path, len(p))
		copy(x, p)
	}

	x = append(x, strings.ToUpper(s))
	return x
}
