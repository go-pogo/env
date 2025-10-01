// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"io"
	"reflect"

	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/rawconv"
)

const ErrStructPointerExpected errors.Msg = "expected a non-nil pointer to a struct"

var unmarshaler rawconv.Unmarshaler

// Unmarshaler is the interface implemented by types that can unmarshal a
// textual representation of themselves.
// It is similar to [encoding.TextUnmarshaler].
type Unmarshaler interface {
	UnmarshalEnv([]byte) error
}

// Unmarshal parses the env formatted data and stores the result in the value
// pointed to by v. If v is nil or not a pointer, Unmarshal returns an
// [ErrStructPointerExpected] error.
func Unmarshal(data []byte, v any) error {
	return NewReaderDecoder(bytes.NewReader(data)).Decode(v)
}

type DecodeOptions struct {
	// ReplaceVars
	ReplaceVars bool
}

// A Decoder looks up environment variables while decoding them into a struct.
type Decoder struct {
	DecodeOptions
	TagOptions

	lookupper Lookupper
}

const panicNilLookupper = "env.Decoder: Lookupper must not be nil"

// NewDecoder returns a new [Decoder] which looks up environment variables from
// the provided [Lookupper](s). When a [Chain] is provided it must not be empty.
func NewDecoder(src ...Lookupper) *Decoder {
	l, chained := chain(src...)
	if !chained && l == nil {
		panic(panicNilLookupper)
	} else if c, ok := l.(chainLookupper); ok && len(c) == 0 {
		panic(panicNilLookupper)
	}

	return (&Decoder{lookupper: l}).
		WithOptions(DecodeOptions{ReplaceVars: true}).
		WithTagOptions(envtag.DefaultOptions())
}

// NewReaderDecoder returns a new [Decoder] similar to calling [NewDecoder] with
// [NewReader] as argument.
func NewReaderDecoder(r io.Reader) *Decoder {
	return (&Decoder{lookupper: NewReader(r)}).
		WithOptions(DecodeOptions{ReplaceVars: true}).
		WithTagOptions(envtag.DefaultOptions())
}

// Strict sets the StrictTags option to true.
func (d *Decoder) Strict() *Decoder {
	d.StrictTags = true
	return d
}

// WithLookupper sets the internal [Lookupper] to l.
func (d *Decoder) WithLookupper(l Lookupper) *Decoder {
	if l == nil {
		panic(panicNilLookupper)
	}

	d.lookupper = l
	return d
}

// WithOptions sets DecodeOptions to the provided [DecodeOptions] opts.
func (d *Decoder) WithOptions(opts DecodeOptions) *Decoder {
	d.DecodeOptions = opts
	return d
}

// WithTagOptions sets TagOptions to the provided [TagOptions] opts.
func (d *Decoder) WithTagOptions(opts TagOptions) *Decoder {
	d.TagOptions = opts
	return d
}

func (d *Decoder) Decode(v any) error {
	if d.lookupper == nil {
		panic(panicNilLookupper)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New(ErrStructPointerExpected)
	}
	if underlyingKind(rv.Elem().Type()) != reflect.Struct {
		return errors.New(ErrStructPointerExpected)
	}

	if d.ReplaceVars {
		if _, ok := d.lookupper.(*Replacer); !ok {
			d.lookupper = NewReplacer(d.lookupper)
		}
	}

	return (&traverser{
		TagOptions:  d.TagOptions,
		isTypeKnown: typeKnownByUnmarshaler,
		handleField: d.decodeField,
	}).start(rv)
}

func (d *Decoder) decodeField(rv reflect.Value, tag envtag.Tag) error {
	val, err := d.lookupper.Lookup(tag.Name)
	if err != nil && !IsNotFound(err) {
		return err
	}
	if val.String() == "" {
		if tag.Default == "" {
			return nil
		}
		val = Value(tag.Default)
	}

	return unmarshaler.Unmarshal(val, rv)
}
