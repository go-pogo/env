// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/rawconv"
	"github.com/go-pogo/writing"
)

const ErrStructExpected errors.Msg = "expected a struct type"

var marshaler rawconv.Marshaler

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
	// ExportPrefix adds an export prefix to each relevant line.
	ExportPrefix bool
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
	var enc Encoder
	return enc.WithTagOptions(envtag.DefaultOptions()).WithWriter(w)
}

// WithOptions changes the internal [EncodeOptions] to opts.
func (e *Encoder) WithOptions(opts EncodeOptions) *Encoder {
	e.EncodeOptions = opts
	return e
}

// WithTagOptions changes the internal [TagOptions] to opts.
func (e *Encoder) WithTagOptions(opts TagOptions) *Encoder {
	e.TagOptions = opts
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
func (e *Encoder) Encode(v any) error {
	switch src := v.(type) {
	case Map:
		for key, val := range src {
			e.writeKeyValueQuoted(key, val.String())
		}
		return nil

	case map[string]Value:
		for key, val := range src {
			e.writeKeyValueQuoted(key, val.String())
		}
		return nil

	case map[fmt.Stringer]Value:
		for key, val := range src {
			e.writeKeyValueQuoted(key.String(), val.String())
		}
		return nil

	case []NamedValue:
		for _, nv := range src {
			e.writeKeyValueQuoted(nv.Name, nv.Value.String())
		}
		return nil

	case []envtag.Tag:
		for _, t := range src {
			e.writeKeyValueQuoted(t.Name, t.Default)
		}
		return nil

	default:
		rv := reflect.ValueOf(src)
		if underlyingKind(rv.Type()) != reflect.Struct {
			return errors.New(ErrStructExpected)
		}

		return (&traverser{
			TagOptions:  e.TagOptions,
			isTypeKnown: typeKnownByMarshaler,
			handleField: e.encodeField,
		}).start(rv)
	}
}

func (e *Encoder) encodeField(rv reflect.Value, tag envtag.Tag) (err error) {
	val := tag.Default
	if e.TakeValues && !rv.IsZero() {
		val, err = marshalAndQuote(rv)
		if err != nil {
			return err
		}
	} else if !e.TakeValues && val == "" {
		val, err = marshalAndQuote(reflect.New(rv.Type()).Elem())
		if err != nil {
			return err
		}
	}

	e.writeKeyValue(tag.Name, val)
	return nil
}

func (e *Encoder) writeKeyValue(key, val string) {
	if e.ExportPrefix {
		_, _ = e.w.WriteString("export ")
	}

	_, _ = e.w.WriteString(key)
	_, _ = e.w.WriteString("=")
	_, _ = e.w.WriteString(val)
	_, _ = e.w.WriteString("\n")
}

func (e *Encoder) writeKeyValueQuoted(key, val string) {
	e.writeKeyValue(key, quote(val))
}

func marshalAndQuote(rv reflect.Value) (string, error) {
	v, err := marshaler.Marshal(rv)
	if err != nil {
		return "", err
	}
	return quote(v.String()), nil
}

func quote(str string) string {
	if str == "" {
		return str
	}

	isq := strings.IndexRune(str, '\'')
	idq := strings.IndexRune(str, '"')
	if isq == -1 && idq == -1 {
		return str
	}

	quot := "\""
	if isq == -1 && idq >= 0 {
		quot = "'"
	} else if isq >= 0 && idq >= 0 {
		str = strings.ReplaceAll(str, quot, "\\"+quot)
	}

	return quot + str + quot
}
