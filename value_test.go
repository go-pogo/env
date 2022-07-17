// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"strconv"
	"testing"
	"time"

	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	types := map[string][]string{
		"empty":  {""},
		"bool":   {"true", "false"},
		"int":    {"100", "+33", "-349"},
		"float":  {"1.1", "0.59999", "22.564856"},
		"time":   {"10s", "2h"},
		"string": {"some value", "another string"},
	}

	tests := map[string][2]func(s string) (interface{}, error){
		"String": {
			func(s string) (interface{}, error) { return Value(s).String(), nil },
			func(s string) (interface{}, error) { return s, nil },
		},
		"TryBool": {
			func(s string) (interface{}, error) { return Value(s).TryBool() },
			func(s string) (interface{}, error) { return strconv.ParseBool(s) },
		},
		"TryInt": {
			func(s string) (interface{}, error) { return Value(s).TryInt() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, strconv.IntSize)
				return int(i), err
			},
		},
		"TryInt8": {
			func(s string) (interface{}, error) { return Value(s).TryInt8() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, 8)
				return int8(i), err
			},
		},
		"TryInt16": {
			func(s string) (interface{}, error) { return Value(s).TryInt16() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, 16)
				return int16(i), err
			},
		},
		"TryInt32": {
			func(s string) (interface{}, error) { return Value(s).TryInt32() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseInt(s, 0, 32)
				return int32(i), err
			},
		},
		"TryInt64": {
			func(s string) (interface{}, error) { return Value(s).TryInt64() },
			func(s string) (interface{}, error) { return strconv.ParseInt(s, 0, 64) },
		},
		"Uint": {
			func(s string) (interface{}, error) { return Value(s).Uint() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, strconv.IntSize)
				return uint(i), err
			},
		},
		"Uint8": {
			func(s string) (interface{}, error) { return Value(s).Uint8() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, 8)
				return uint8(i), err
			},
		},
		"Uint16": {
			func(s string) (interface{}, error) { return Value(s).Uint16() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, 16)
				return uint16(i), err
			},
		},
		"Uint32": {
			func(s string) (interface{}, error) { return Value(s).Uint32() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseUint(s, 0, 32)
				return uint32(i), err
			},
		},
		"Uint64": {
			func(s string) (interface{}, error) { return Value(s).Uint64() },
			func(s string) (interface{}, error) { return strconv.ParseUint(s, 0, 64) },
		},
		"Float32": {
			func(s string) (interface{}, error) { return Value(s).Float32() },
			func(s string) (interface{}, error) {
				i, err := strconv.ParseFloat(s, 32)
				return float32(i), err
			},
		},
		"Float64": {
			func(s string) (interface{}, error) { return Value(s).Float64() },
			func(s string) (interface{}, error) { return strconv.ParseFloat(s, 64) },
		},
		"Duration": {
			func(s string) (interface{}, error) { return Value(s).Duration() },
			func(s string) (interface{}, error) { return time.ParseDuration(s) },
		},
	}

	for name, tcFn := range tests {
		t.Run(name, func(t *testing.T) {
			for typ, inputs := range types {
				for _, input := range inputs {
					t.Run(typ, func(t *testing.T) {
						haveVal, haveErr := tcFn[0](input)
						wantVal, wantErr := tcFn[1](input)

						assert.Exactly(t, wantVal, haveVal)
						assert.Exactly(t, wantErr, errors.Unwrap(haveErr))

						if wantErr != nil {
							haveKind := errors.GetKind(haveErr)
							assert.True(t, haveKind == ParseError || haveKind == ValidationError, "Kind should match ParseError or ValidationError")
						}
					})
				}
			}
		})
	}
}
