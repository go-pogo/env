// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUnderlyingKind(t *testing.T) {
	tests := map[string]struct {
		value    interface{}
		wantKind reflect.Kind
	}{
		"string": {
			value:    "",
			wantKind: reflect.String,
		},
		"*string": {
			value:    (*string)(nil),
			wantKind: reflect.String,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.wantKind, underlyingKind(reflect.TypeOf(tc.value)))
		})
	}
}

func BenchmarkUnderlyingKind(b *testing.B) {
	types := []reflect.Type{
		reflect.TypeOf(""),
		reflect.TypeOf((*string)(nil)),
		reflect.TypeOf((*****string)(nil)),
	}
	for _, typ := range types {
		b.Run("loop_"+typ.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				loopUnderlyingKind(typ)
			}
		})
		b.Run("recursive_"+typ.String(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				recursiveUnderlyingKind(typ)
			}
		})
	}
}

//go:noinline
func loopUnderlyingKind(rt reflect.Type) reflect.Kind {
	k := rt.Kind()
	for k == reflect.Ptr {
		rt = rt.Elem()
		k = rt.Kind()
	}
	return k
}

//go:noinline
func recursiveUnderlyingKind(rt reflect.Type) reflect.Kind {
	if k := rt.Kind(); k != reflect.Ptr {
		return k
	}
	return recursiveUnderlyingKind(rt.Elem())
}
