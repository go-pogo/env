// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bytes"
	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/rawconv"
	"io"
	"reflect"
)

const ErrStructPointerExpected errors.Msg = "expected a non-nil pointer to a struct"

var unmarshaler rawconv.Unmarshaler

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

type Decoder struct {
	envtag.Options
	lookupper Lookupper

	// ReplaceVars
	ReplaceVars bool
}

const panicNilLookupper = "env.Decoder: Lookupper must not be nil"

// NewDecoder returns a new Decoder that looks up environment variables from
// any Lookupper.
func NewDecoder(src ...Lookupper) *Decoder {
	l, chained := chain(src...)
	if !chained && l == nil {
		panic(panicNilLookupper)
	} else if c, ok := l.(chainLookupper); ok && len(c) == 0 {
		panic(panicNilLookupper)
	}

	return &Decoder{
		lookupper:   l,
		Options:     envtag.DefaultOptions(),
		ReplaceVars: true,
	}
}

// NewReaderDecoder returns a new Decoder similar to calling NewDecoder with
// NewReader as argument.
func NewReaderDecoder(r io.Reader) *Decoder {
	return &Decoder{
		lookupper:   NewReader(r),
		Options:     envtag.DefaultOptions(),
		ReplaceVars: true,
	}
}

// Strict sets the StrictTags option to true.
func (d *Decoder) Strict() *Decoder {
	d.Options.StrictTags = true
	return d
}

// WithLookupper changes the Decoder's internal Lookupper to l.
func (d *Decoder) WithLookupper(l Lookupper) *Decoder {
	if l == nil {
		panic(panicNilLookupper)
	}

	d.lookupper = l
	return d
}

func (d *Decoder) WithOptions(opts envtag.Options) *Decoder {
	d.Options = opts
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
	if underlyingKind(rv.Type()) != reflect.Struct {
		return errors.New(ErrStructPointerExpected)
	}

	if d.ReplaceVars {
		if _, ok := d.lookupper.(*Replacer); !ok {
			d.lookupper = NewReplacer(d.lookupper)
		}
	}

	return (&traverser{
		Options:     d.Options,
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
