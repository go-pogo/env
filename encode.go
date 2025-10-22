// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/writing"
)

const ErrStructExpected errors.Msg = "expected a struct type"

// Marshaler is the interface implemented by types that can marshal themselves
// into valid env values.
type Marshaler interface {
	MarshalEnv() ([]byte, error)
}

// Marshal returns v encoded in env format.
func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type EncodeOptions struct {
	// TakeValues takes the values from the struct field.
	TakeValues bool
	// ExportPrefix adds an "export" prefix to each relevant line.
	// Deprecated: Use [FormatShellExport] [Formatter] instead.
	ExportPrefix bool

	Formatter Formatter
}

// An Encoder writes env values to an output stream.
type Encoder struct {
	EncodeOptions
	TagOptions

	w writing.StringWriter
}

const panicNilWriter = "env.Encoder: io.Writer must not be nil"

// NewEncoder returns a new [Encoder] which writes to w.
func NewEncoder(w io.Writer) *Encoder {
	if w == nil {
		panic(panicNilWriter)
	}

	return &Encoder{
		EncodeOptions: EncodeOptions{Formatter: Format},
		TagOptions:    envtag.DefaultOptions(),
		w:             writing.ToStringWriter(w),
	}
}

// WithOptions sets EncodeOptions to the provided [EncodeOptions] opts.
func (e *Encoder) WithOptions(opts EncodeOptions) *Encoder {
	e.EncodeOptions = opts
	return e
}

// WithTagOptions sets TagOptions to the provided [TagOptions] opts.
func (e *Encoder) WithTagOptions(opts TagOptions) *Encoder {
	e.TagOptions = opts
	return e
}

// WithFormatter sets Formatter to the provided [Formatter] p.
func (e *Encoder) WithFormatter(p Formatter) *Encoder {
	e.Formatter = p
	return e
}

// WithWriter changes the internal [io.Writer] to w.
func (e *Encoder) WithWriter(w io.Writer) *Encoder {
	if w == nil {
		panic(panicNilWriter)
	}
	e.w = writing.ToStringWriter(w)
	return e
}

// Encode writes the env format encoding of v to the underlying [io.Writer].
// Supported types of v are:
//   - [Map]
//   - map[string][Value]
//   - map[[fmt.Stringer]][Value]
//   - [][NamedValue]
//   - [][envtag.Tag]
//   - any struct type the rawconv package can handle
func (e *Encoder) Encode(v any) (err error) {
	if e.Formatter == nil {
		//goland:noinspection GoDeprecation
		if e.ExportPrefix {
			e.Formatter = FormatShellExport
		} else {
			e.Formatter = Format
		}
	}

	switch src := v.(type) {
	case Map:
		for key, val := range src {
			if err = e.print(key, val); err != nil {
				return err
			}
		}
		return nil

	case map[string]Value:
		for key, val := range src {
			if err = e.print(key, val); err != nil {
				return err
			}
		}
		return nil

	case map[fmt.Stringer]Value:
		for key, val := range src {
			if err = e.print(key.String(), val); err != nil {
				return err
			}
		}
		return nil

	case []NamedValue:
		for _, nv := range src {
			if err = e.print(nv.Name, nv.Value); err != nil {
				return err
			}
		}
		return nil

	case []envtag.Tag:
		for _, t := range src {
			if err = e.print(t.Name, t.Default); err != nil {
				return err
			}
		}
		return nil

	default:
		rv := reflect.ValueOf(src)
		if underlyingKind(rv.Type()) != reflect.Struct {
			return errors.New(ErrStructExpected)
		}

		return (&traverser{
			TagOptions:  e.TagOptions,
			isKnownType: typeKnownByMarshaler,
			handleField: e.encodeField,
		}).start(rv)
	}
}

func (e *Encoder) encodeField(rv reflect.Value, tag envtag.Tag) error {
	if !e.TakeValues && tag.Default == "" {
		return e.print(tag.Name, reflect.New(rv.Type()).Elem())
	}
	if tag.Default != "" && (!e.TakeValues || (e.TakeValues && rv.IsZero())) {
		var err error
		if rv, err = defaultValue(rv.Type(), tag.DefaultValue()); err != nil {
			return err
		}
	}

	return e.print(tag.Name, rv)
}

func (e *Encoder) print(name string, val any) error {
	str, err := e.Formatter(name, val)
	if err != nil {
		return err
	}
	if _, err = e.w.WriteString(str + "\n"); err != nil {
		return err
	}

	return nil
}

func defaultValue(t reflect.Type, v Value) (reflect.Value, error) {
	ptr := reflect.New(t)
	err := unmarshaler.Unmarshal(v, ptr)
	return ptr.Elem(), err
}
