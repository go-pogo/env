// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/parseval"
	"io"
	"reflect"
)

const (
	ErrStructPointerExpected errors.Msg = "expected a non-nil pointer to a struct"
)

var unmarshaler parseval.Unmarshaler

func init() {
	unmarshaler.Register(
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

// Unmarshal parses the env formatted data and stores the result in the value
// pointed to by v. If v is nil or not a pointer, Unmarshal returns an
// ErrStructPointerExpected error.
func Unmarshal(data []byte, v any) error {
	return NewDecoder(NewReader(bytes.NewReader(data))).Decode(v)
}

var _ Lookupper = new(Decoder)

type Decoder struct {
	envtag.Options
	Lookupper
	r *Replacer

	// ReplaceVars
	ReplaceVars bool
}

// NewDecoder returns a new Decoder that looks up environment variables from
// any Lookupper.
//
//	dec := NewDecoder(NewReader(r))
func NewDecoder(src ...Lookupper) *Decoder {
	return &Decoder{
		Options:   envtag.DefaultOptions(),
		Lookupper: Chain(src...),
	}
}

// NewReaderDecoder returns a new Decoder similar to calling NewDecoder with
// NewReader.
func NewReaderDecoder(r io.Reader) *Decoder {
	return &Decoder{
		Options:   envtag.DefaultOptions(),
		Lookupper: NewReader(r),
	}
}

// WithLookupper changes the Decoder's internal Lookupper to l.
func (d *Decoder) WithLookupper(l Lookupper) *Decoder {
	d.Lookupper = l
	d.r = nil
	return d
}

func (d *Decoder) WithOptions(opts envtag.Options) *Decoder {
	d.Options = opts
	return d
}

const panicNilLookupper = "env.Decoder: Lookupper must not be nil"

func (d *Decoder) Decode(v any) error {
	if d.Lookupper == nil {
		panic(panicNilLookupper)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New(ErrStructPointerExpected)
	}
	if underlyingKind(rv.Type()) != reflect.Struct {
		return errors.New(ErrStructPointerExpected)
	}

	return newTraverser(d.Options, d.decodeField).start(rv)
}

func (d *Decoder) decodeField(rv reflect.Value, tag envtag.Tag) error {
	val, err := d.Lookup(tag.Name)
	if err != nil {
		if !IsNotFound(err) {
			return err
		}
		if tag.Default == "" {
			return nil
		}
		val = Value(tag.Default)
	}

	return unmarshaler.Unmarshal(val, rv)
}

func (d *Decoder) Lookup(key string) (Value, error) {
	if d.Lookupper == nil {
		panic(panicNilLookupper)
	}
	if !d.ReplaceVars {
		return d.Lookupper.Lookup(key)
	}
	if d.r == nil {
		d.r = NewReplacer(d.Lookupper)
	}
	return d.r.Lookup(key)
}
