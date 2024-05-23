// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envfile

import (
	"github.com/go-pogo/env"
	"github.com/go-pogo/errors"
	"io"
	"os"
)

var _ io.Closer = (*Encoder)(nil)

type encoder = env.Encoder

// Encoder embeds an env.Encoder and sets its target io.Writer to an os.File.
type Encoder struct {
	*encoder
	file *os.File
}

// NewEncoder returns a new Encoder which writes the encoded values to the
// provided os.File f.
func NewEncoder(f *os.File) *Encoder {
	if f == nil {
		panic(panicNilFile)
	}
	return &Encoder{
		encoder: env.NewEncoder(f),
		file:    f,
	}
}

// Create returns a new Encoder which creates the file filename and uses it to
// write the encoded values to.
func Create(filename string) (*Encoder, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewEncoder(f), nil
}

// Write creates the file filename using Create, writes the encoded value v
// to it and closes its internal file handle.
func Write(filename string, v any) (err error) {
	enc, err := Create(filename)
	if err != nil {
		return err
	}

	defer errors.AppendFunc(&err, enc.Close)
	if err = enc.Encode(v); err != nil {
		return err
	}
	return nil
}

// Close closes its internal os.File.
func (fe *Encoder) Close() error {
	return errors.WithStack(fe.file.Close())
}
