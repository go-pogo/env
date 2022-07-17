// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"encoding"
	"reflect"
	"strconv"
	"time"

	"github.com/go-pogo/errors"
)

const ValidationError errors.Kind = "validation error"

// Value is a textual representation of a value which is able to cast itself to
// any of the supported types using its corresponding method.
//
//  boolVal, err := env.Value("true").TryBool()
type Value string

// Empty indicates if Value is an empty string.
func (v Value) Empty() bool { return string(v) == "" }

func (v Value) GoString() string { return `env.Value("` + string(v) + `")` }

// String returns Value as a string.
func (v Value) String() string { return string(v) }

func (v Value) Bool() bool {
	x, _ := strconv.ParseBool(string(v))
	return x
}

// TryBool tries to parse Value as a bool with strconv.ParseBool.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns an error.
func (v Value) TryBool() (bool, error) {
	x, err := strconv.ParseBool(string(v))
	return x, errors.WithKind(err, errKind(err))
}

// TryInt tries to parse Value as an int with strconv.ParseInt.
func (v Value) TryInt() (int, error) {
	x, err := v.tryIntSize(strconv.IntSize)
	return int(x), err
}

// TryInt8 tries to parse Value as an int8 with strconv.ParseInt.
func (v Value) TryInt8() (int8, error) {
	x, err := v.tryIntSize(8)
	return int8(x), err
}

// TryInt16 tries to parse Value as an int16 with strconv.ParseInt.
func (v Value) TryInt16() (int16, error) {
	x, err := v.tryIntSize(16)
	return int16(x), err
}

// TryInt32 tries to parse Value as an int32 with strconv.ParseInt.
func (v Value) TryInt32() (int32, error) {
	x, err := v.tryIntSize(32)
	return int32(x), err
}

// TryInt64 tries to parse Value as an int64 with strconv.ParseInt.
func (v Value) TryInt64() (int64, error) {
	return v.tryIntSize(64)
}

func (v Value) tryIntSize(bitSize int) (int64, error) {
	x, err := strconv.ParseInt(string(v), 0, bitSize)
	return x, errors.WithKind(err, errKind(err))
}

// Uint tries to parse Value as an uint with strconv.ParseUint.
func (v Value) Uint() (uint, error) {
	x, err := v.uintSize(strconv.IntSize)
	return uint(x), err
}

// Uint8 tries to parse Value as an uint8 with strconv.ParseUint.
func (v Value) Uint8() (uint8, error) {
	x, err := v.uintSize(8)
	return uint8(x), err
}

// Uint16 tries to parse Value as an uint16 with strconv.ParseUint.
func (v Value) Uint16() (uint16, error) {
	x, err := v.uintSize(16)
	return uint16(x), err
}

// Uint32 tries to parse Value as an uint32 with strconv.ParseUint.
func (v Value) Uint32() (uint32, error) {
	x, err := v.uintSize(32)
	return uint32(x), err
}

// Uint64 tries to parse Value as an uint64 with strconv.ParseUint.
func (v Value) Uint64() (uint64, error) {
	return v.uintSize(64)
}

func (v Value) uintSize(bitSize int) (uint64, error) {
	x, err := strconv.ParseUint(string(v), 0, bitSize)
	return x, errors.WithKind(err, errKind(err))
}

// Float32 tries to parse Value as a float32 with strconv.ParseFloat.
func (v Value) Float32() (float32, error) {
	x, err := v.floatSize(32)
	return float32(x), err
}

// Float64 tries to parse Value as a float64 with strconv.ParseFloat.
func (v Value) Float64() (float64, error) {
	return v.floatSize(64)
}

func (v Value) floatSize(bitSize int) (float64, error) {
	x, err := strconv.ParseFloat(string(v), bitSize)
	return x, errors.WithKind(err, errKind(err))
}

// Complex64 tries to parse Value as a complex64 with strconv.ParseComplex.
func (v Value) Complex64() (complex64, error) {
	x, err := v.complexSize(64)
	return complex64(x), err
}

// Complex128 tries to parse Value as a complex128 with strconv.ParseComplex.
func (v Value) Complex128() (complex128, error) {
	return v.complexSize(128)
}

func (v Value) complexSize(bitSize int) (complex128, error) {
	x, err := strconv.ParseComplex(string(v), bitSize)
	return x, errors.WithKind(err, errKind(err))
}

// Duration tries to parse Value as time.Duration with time.ParseDuration.
func (v Value) Duration() (time.Duration, error) {
	x, err := time.ParseDuration(string(v))
	if err != nil {
		err = errors.WithKind(err, ParseError)
	}
	return x, err
}

var (
	unmarshalerType   = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	unmarshalTextType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	timeDurationType  = reflect.TypeOf(time.Nanosecond)
)

func (v Value) ReflectAssign(dest reflect.Value) error {
	typ := dest.Type()
	if typ.Implements(unmarshalerType) {
		err := dest.Interface().(Unmarshaler).UnmarshalEnv([]byte(v))
		return errors.WithKind(err, errKind(err))
	}
	if typ.Implements(unmarshalTextType) {
		err := dest.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(v))
		return errors.WithKind(err, errKind(err))
	}
	if v.Empty() {
		return nil
	}
	if typ == timeDurationType {
		if x, err := v.Duration(); err != nil {
			return err
		} else {
			dest.Set(reflect.ValueOf(x))
			return nil
		}
	}

	switch dest.Kind() {
	case reflect.String:
		dest.SetString(v.String())
		return nil

	case reflect.Bool:
		x, err := v.TryBool()
		dest.SetBool(x)
		return err

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := v.tryIntSize(typ.Bits())
		dest.SetInt(x)
		return err

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := v.uintSize(typ.Bits())
		dest.SetUint(x)
		return err

	case reflect.Float32, reflect.Float64:
		x, err := v.floatSize(typ.Bits())
		dest.SetFloat(x)
		return err

	case reflect.Complex64, reflect.Complex128:
		x, err := v.complexSize(typ.Bits())
		dest.SetComplex(x)
		return err

	default:
		return &UnsupportedError{Type: typ}
	}
}

type UnsupportedError struct {
	Type reflect.Type
}

func (e *UnsupportedError) Error() string {
	return "type `" + e.Type.String() + "` is unsupported"
}

func errKind(err error) errors.Kind {
	if ne, ok := err.(*strconv.NumError); ok {
		if ne.Err == strconv.ErrRange {
			return ValidationError
		} else {
			return ParseError
		}
	}
	return errors.UnknownKind
}
