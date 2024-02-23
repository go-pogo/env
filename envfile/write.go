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

type Encoder struct {
	*encoder
	file *os.File
}

func NewEncoder(f *os.File) *Encoder {
	if f == nil {
		panic(panicNilFile)
	}
	return &Encoder{
		encoder: env.NewEncoder(f),
		file:    f,
	}
}

func Create(filename string) (*Encoder, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewEncoder(f), nil
}

func Write(filename string, v any) error {
	enc, err := Create(filename)
	if err != nil {
		return err
	}
	return enc.Encode(v)
}

// Close closes the underlying os.File.
func (fe *Encoder) Close() error {
	return errors.WithStack(fe.file.Close())
}
