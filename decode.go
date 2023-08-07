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

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(NewReader(bytes.NewBuffer(data))).Decode(v)
}

var _ Lookupper = new(Decoder)

type Decoder struct {
	envtag.Options
	Lookupper
	traverser

	// ReplaceVars
	ReplaceVars bool
}

// NewDecoder returns a new Decoder that looks up environment variables from
// any Lookupper.
//
//	dec := NewDecoder(NewReader(r))
func NewDecoder(src ...Lookupper) *Decoder {
	var d Decoder
	d.init(Chain(src...))
	return &d
}

// NewReaderDecoder returns a new Decoder similar to calling NewDecoder with
// NewReader. It looks up environment variables from the new Reader.
func NewReaderDecoder(r io.Reader) *Decoder {
	var d Decoder
	d.init(NewReader(r))
	return &d
}

func (d *Decoder) init(l Lookupper) {
	d.Lookupper = l
	d.ReplaceVars = true
	d.Options.Defaults()
	d.traverser.Options = &d.Options
	d.traverser.HandleField = d.decodeField
}

func (d *Decoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New(ErrStructPointerExpected)
	}
	if underlyingKind(rv.Type()) != reflect.Struct {
		return errors.New(ErrStructPointerExpected)
	}

	return d.traverser.traverse(rv, "")
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
	if !d.ReplaceVars {
		return d.Lookupper.Lookup(key)
	}

	// todo: check + replace vars in value
	return d.Lookupper.Lookup(key)
}
