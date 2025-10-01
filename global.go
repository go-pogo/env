// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"reflect"

	"github.com/go-pogo/rawconv"
)

var (
	marshaler   rawconv.Marshaler
	unmarshaler rawconv.Unmarshaler
)

func init() {
	marshaler.Register(
		reflect.TypeOf((*Marshaler)(nil)).Elem(),
		func(v any) (string, error) {
			b, err := v.(Marshaler).MarshalEnv()
			if err != nil {
				return "", err
			}
			return string(b), err
		},
	)

	unmarshaler.Register(
		reflect.TypeOf((*Unmarshaler)(nil)).Elem(),
		func(val rawconv.Value, dest any) error {
			return dest.(Unmarshaler).UnmarshalEnv(val.Bytes())
		},
	)
}

// GetMarshalFunc returns the globally registered [rawconv.MarshalFunc] for
// [reflect.Type] typ or nil if there is none registered.
func GetMarshalFunc(typ reflect.Type) rawconv.MarshalFunc {
	return marshaler.Func(typ)
}

// GetUnmarshalFunc returns the globally registered [rawconv.UnmarshalFunc] for
// [reflect.Type] typ or nil if there is none registered.
func GetUnmarshalFunc(typ reflect.Type) rawconv.UnmarshalFunc {
	return unmarshaler.Func(typ)
}
