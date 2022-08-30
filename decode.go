// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"io"
	"reflect"
	"strings"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/parseval"
)

// Unmarshaler is the interface implemented by types that can unmarshal a
// textual representation of themselves.
// It is similar to encoding.TextUnmarshaler.
type Unmarshaler interface {
	UnmarshalEnv([]byte) error
}

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewBuffer(data), UnmarshalOpts).Decode(v)
}

// Decode reads from io.Reader and decodes its relevant content to v.
// It matches the configure.DecoderFunc signature.
func Decode(r io.Reader, v interface{}) error {
	return NewDefaultDecoder(r).Decode(v)
}

const (
	TagsOnly    Option = 1 << iota // ignore fields that do not have an `env` tag
	StripExport                    // strip "export " from start of line
	ReplaceVars                    // replace ${} vars

	UnmarshalOpts  = StripExport | ReplaceVars
	DefaultOptions = TagsOnly | UnmarshalOpts

	ErrStructPointerExpected errors.Msg = "expected a pointer to a struct"

	Tag      = "env"
	ignore   = "-"
	noPrefix = "noprefix"
)

type Option uint8

func (o Option) has(f Option) bool { return o&f != 0 }

type Decoder struct {
	scanner  Scanner
	fallback Lookupper
	parser   *parseval.Parser
	err      error
	found    Map
	opts     Option
}

// NewDecoder returns a new Decoder that scans io.Reader r for environment
// variables and parses them.
func NewDecoder(r io.Reader, opts Option) *Decoder {
	return (&Decoder{
		parser: parseval.NewParser(
			reflect.TypeOf((*Unmarshaler)(nil)).Elem(),
			func(v parseval.Value, u interface{}) error {
				return u.(Unmarshaler).UnmarshalEnv([]byte(v))
			},
		),
		opts: opts,
	}).Reset(r)
}

// NewDefaultDecoder uses NewDecoder to create a *Decoder with DefaultOptions
// and LookupEnv as Fallback.
func NewDefaultDecoder(r io.Reader) *Decoder {
	return setDefaultFallback(NewDecoder(r, DefaultOptions))
}

func setDefaultFallback(dec *Decoder) *Decoder {
	dec.fallback = LookupperFunc(LookupEnv)
	return dec
}

const panicSelfAsFallback = "cannot set Decoder as a fallback of itself"

// Fallback Lookupper which is called when a key cannot be found within the
// internal io.Reader that's used to construct the Decoder.
func (d *Decoder) Fallback(fallback Lookupper) {
	if d == fallback {
		panic(panicSelfAsFallback)
	}
	d.fallback = fallback
}

func (d *Decoder) Reset(r io.Reader) *Decoder {
	d.scanner = NewScanner(r)
	d.found = make(Map, 4)
	return d
}

// Options return the set Option flags.
func (d *Decoder) Options() Option { return d.opts }

// Err returns an error that might occur during scanning when directly using
// Lookup.
func (d *Decoder) Err() error {
	if d.err == nil {
		return nil
	}

	err := d.err
	d.err = nil
	return err
}

// Lookup retrieves the Value of the environment variable named by key.
// If the variable is present the value (which may be empty) is returned and
// the boolean is true. Otherwise, the returned value will be empty and the
// boolean will be false.
func (d *Decoder) Lookup(key string) (Value, bool) {
	v, ok, err := d.lookup(key, true)
	if err != nil {
		errors.Append(&d.err, err)
	}
	return v, ok
}

type LookupError struct {
	Err error
	Key string
}

func (e *LookupError) Unwrap() error { return e.Err }

func (e *LookupError) Error() string { return "error while looking up `" + e.Key + "`" }

func (d *Decoder) lookup(lookup string, fallback bool) (Value, bool, error) {
	if val, ok := d.found[lookup]; ok {
		return val, true, nil
	}

	// defer errors.CatchPanic(err)
	for d.scanner.Scan() {
		key, val, err := parseAndStore(d.found, d.scanner.Text(), d.opts.has(StripExport))
		if err != nil {
			return "", false, errors.WithStack(&LookupError{
				Err: err,
				Key: lookup,
			})
		}

		// found the key we were looking for
		// no need to continue scanning, for now...
		if key == lookup {
			return val, true, nil
		}
	}
	if fallback && d.fallback != nil {
		if val, ok := d.fallback.Lookup(lookup); ok {
			d.found[lookup] = val
			return val, true, nil
		}
	}
	return "", false, nil
}

func (d *Decoder) scanAll() error {
	return scanAll(d.scanner, d.found, d.opts.has(StripExport))
}

// Map returns a Map of all found environment variables.
func (d *Decoder) Map() (Map, error) {
	if err := d.scanAll(); err != nil {
		return nil, err
	}

	clone := make(Map, len(d.found))
	clone.MergeValues(d.found)
	return clone, nil
}

func (d *Decoder) Decode(v interface{}) error {
	if m, ok := v.(Map); ok {
		if err := d.scanAll(); err != nil {
			return err
		}

		m.MergeValues(d.found)
		return nil
	}

	if v == nil || reflect.TypeOf(v).Kind() != reflect.Ptr {
		return ErrStructPointerExpected
	}

	val := reflect.ValueOf(v).Elem()
	if val.Kind() != reflect.Struct {
		return ErrStructPointerExpected
	}

	return d.traverseStruct(val.Type(), val, "")
}

func (d *Decoder) DecodeField(field reflect.StructField, val reflect.Value, lookup string) error {
	tag, found := field.Tag.Lookup(Tag)
	if tag == ignore || (!found && d.opts.has(TagsOnly)) {
		return nil
	}
	if tag != "" {
		lookup = tag
	}
	if lookup == "" {
		lookup = field.Name
	}

	v, ok, err := d.lookup(lookup, true)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	return d.parser.Parse(v, val)
}

func (d *Decoder) traverseStruct(pt reflect.Type, pv reflect.Value, p string) error {
	if len(p) > 0 {
		p += "_"
	}

	for i := 0; i < pv.NumField(); i++ {
		field, val := pt.Field(i), pv.Field(i)
		path := p + strings.ToUpper(field.Name)

		switch field.Type.Kind() {
		case reflect.Ptr:
			elem := field.Type.Elem()
			if !val.CanSet() || elem.Kind() != reflect.Struct {
				continue
			}

			val.Set(reflect.New(elem))
			if err := d.traverseStruct(field.Type, val.Elem(), path); err != nil {
				return err
			}

		case reflect.Struct:
			if err := d.traverseStruct(field.Type, val, path); err != nil {
				return err
			}

		default:
			if err := d.DecodeField(field, val, path); err != nil {
				return err
			}
		}
	}
	return nil
}
