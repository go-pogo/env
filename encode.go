// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/parseval"
	"github.com/go-pogo/writing"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const (
	ErrStructExpected errors.Msg = "expected a struct type"
)

type Marshaler interface {
	MarshalEnv() ([]byte, error)
}

type Encoder struct {
	envtag.Options
	traverser
	w writing.StringWriter

	// TakeValues takes the values from the struct field.
	TakeValues bool
	// ExportPrefix adds an export prefix to each relevant line.
	ExportPrefix bool
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{
		w: writing.ToStringWriter(w),
	}
	e.Options.Defaults()
	e.traverser.Options = &e.Options
	e.traverser.HandleField = e.encodeField
	return e
}

func (e *Encoder) Encode(v interface{}) error {
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
		return e.traverser.traverse(rv, "")
	}
}

func (e *Encoder) encodeField(rv reflect.Value, tag envtag.Tag) error {
	val := tag.Default
	if e.TakeValues {
		if v, err := encodeValue(rv, tag); err != nil {
			return err
		} else if v != "" {
			val = v
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

func encodeValue(rv reflect.Value, tag envtag.Tag) (string, error) {
	for rv.Kind() == reflect.Ptr {
		// todo: check of rv Marshaler of encoding.TextMarshaler implement
		rv = rv.Elem()
	}

	switch underlyingKind(rv.Type()) {
	case reflect.Invalid, reflect.Uintptr, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Slice, reflect.UnsafePointer:
		return "", &parseval.UnsupportedTypeError{Type: rv.Type()}
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmtInt(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmtUint(rv.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmtFloat(rv.Float(), rv.Type().Bits()), nil
	case reflect.Complex64, reflect.Complex128:
		return fmtComplex(rv.Complex(), rv.Type().Bits()), nil
	case reflect.String:
		return quote(rv.String()), nil
	case reflect.Struct:

	}

	return fmt.Sprintf("%v", rv.Interface()), nil
}

func fmtInt(v int64) string { return strconv.FormatInt(v, 10) }

func fmtUint(v uint64) string { return strconv.FormatUint(v, 10) }

func fmtFloat(v float64, bitSize int) string {
	return strconv.FormatFloat(v, 'g', -1, bitSize)
}

func fmtComplex(v complex128, bitSize int) string {
	return strconv.FormatComplex(v, 'g', -1, bitSize)
}
