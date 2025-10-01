// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"reflect"

	"github.com/go-pogo/env/envtag"
)

func typeKnownByUnmarshaler(typ reflect.Type) bool {
	return unmarshaler.Func(typ) != nil
}

func typeKnownByMarshaler(typ reflect.Type) bool {
	return marshaler.Func(typ) != nil
}

type traverser struct {
	TagOptions

	isKnownType func(reflect.Type) bool
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

		opts := t.TagOptions
		opts.StrictTags = opts.StrictTags && !include

		tag, _ := envtag.ParseStructField(opts, field, prefix)
		if tag.ShouldIgnore() {
			continue
		}

		/*if t.isKnownType == nil {
			if handled, err := t.handleField(rv, tag); err != nil {
				return err
			} else if handled {
				continue
			}
		} else*/if kind != reflect.Struct || t.isKnownType(rv.Type()) {
			// when the field is a struct and is a known type (of the global
			// (un)marshaler) it means the Decoder/Encoder can handle the field
			// after which we'll continue with the next field, without further
			// traversing the struct's fields
			if err := t.handleField(rv, tag); err != nil {
				return err
			}
			continue
		}

		var p string
		if tag.Inline || field.Anonymous {
			p = prefix
		} else {
			p = tag.Name
		}

		if err := t.traverse(rv, p, include || tag.Include); err != nil {
			return err
		}
	}
	return nil
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
