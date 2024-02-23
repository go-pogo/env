// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/rawconv"
	"github.com/go-pogo/writing"
	"io"
	"reflect"
	"strings"
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

// An Encoder writes env values to an output stream.
type Encoder struct {
	envtag.Options
	w writing.StringWriter

	// TakeValues takes the values from the struct field.
	TakeValues bool
	// ExportPrefix adds an export prefix to each relevant line.
	ExportPrefix bool
}

const panicNilWriter = "env.Encoder: io.Writer must not be nil"

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	var enc Encoder
	return enc.WithOptions(envtag.DefaultOptions()).WithWriter(w)
}

func (e *Encoder) WithOptions(opts envtag.Options) *Encoder {
	e.Options = opts
	return e
}

func (e *Encoder) WithWriter(writer io.Writer) *Encoder {
	if writer == nil {
		panic(panicNilWriter)
	}
	e.w = writing.ToStringWriter(writer)
	return e
}

// Encode writes the env format encoding of v to the underlying io.Writer.
// Supported types of v are:
//   - Map
//   - map[string]Value
//   - []NamedValue
//   - []envtag.Tag
//   - any valid struct type
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
			Options:     e.Options,
			isTypeKnown: typeKnownByMarshaler,
			handleField: e.encodeField,
		}).start(rv)
	}
}

func (e *Encoder) encodeField(rv reflect.Value, tag envtag.Tag) error {
	val := tag.Default
	if e.TakeValues {
		if v, err := marshaler.Marshal(rv); err != nil {
			return err
		} else if !v.IsEmpty() {
			val = quote(v.String())
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
