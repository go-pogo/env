// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/rawconv"
	"reflect"
)

func init() {
	unmarshaler.Register(
		reflect.TypeOf((*Unmarshaler)(nil)).Elem(),
		func(val rawconv.Value, dest any) error {
			return dest.(Unmarshaler).UnmarshalEnv(val.Bytes())
		},
	)

	marshaler.Register(
		reflect.TypeOf((*Marshaler)(nil)).Elem(),
		func(v any) (string, error) {
			b, err := v.(Marshaler).MarshalEnv()
			if err != nil {
				return "", err
			}
			return string(b), err
		},
	)
}

func typeKnownByUnmarshaler(typ reflect.Type) bool {
	return unmarshaler.Func(typ) != nil
}

func typeKnownByMarshaler(typ reflect.Type) bool {
	return marshaler.Func(typ) != nil
}

type traverser struct {
	envtag.Options

	isTypeKnown func(reflect.Type) bool
	handleField func(reflect.Value, envtag.Tag) error
}

func (t *traverser) start(pv reflect.Value) error {
	return t.traverse(pv, "", false)
}

const panicPtr = "env: ptr values should always be resolved; this is a bug!"

func (t *traverser) traverse(pv reflect.Value, prefix string, include bool) error {
	// todo: dit moet anders, in het geval van encode zou pv.Interface() wss. prima nil kunnen zijn
	pv = indirect(pv)

	pt := pv.Type()
	for i, l := 0, pv.NumField(); i < l; i++ {
		field, rv := pt.Field(i), pv.Field(i)
		kind := underlyingKind(field.Type)
		if kind == reflect.Invalid || kind == reflect.Uintptr || kind == reflect.Chan || kind == reflect.Func || kind == reflect.UnsafePointer {
			// unsupported types
			continue
		} else if kind == reflect.Ptr {
			panic(panicPtr)
		}

		opts := t.Options
		opts.StrictTags = opts.StrictTags && !include

		tag, _ := envtag.ParseStructField(opts, field)
		if tag.ShouldIgnore() {
			continue
		}

		if kind == reflect.Struct && !t.isTypeKnown(rv.Type()) {
			p := prefix
			if tag.NoPrefix {
				// ignore prefix
				p = tag.Name
			} else if !tag.Inline && !field.Anonymous {
				// append tag name to prefix
				p = prefixAppend(prefix, tag.Name)
			}

			if err := t.traverse(rv, p, include || tag.Include); err != nil {
				return err
			} else {
				// no error, continue to next field
				continue
			}
		}

		if !tag.NoPrefix {
			tag.Name = prefixAppend(prefix, tag.Name)
		}
		if err := t.handleField(rv, tag); err != nil {
			return err
		}
	}
	return nil
}

func prefixAppend(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "_" + name
}

// underlyingKind resolves the underlying reflect.Kind of typ.
func underlyingKind(typ reflect.Type) reflect.Kind {
	if k := typ.Kind(); k != reflect.Ptr {
		return k
	}
	return underlyingKind(typ.Elem())
}

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// create a pointer to the type v points to
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}
	return v
}
