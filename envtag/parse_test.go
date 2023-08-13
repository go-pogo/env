// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseTag(t *testing.T) {
	tests := map[string]struct {
		wantTag Tag
		wantErr error
	}{
		"":                 {wantTag: Tag{}},
		"-":                {wantTag: Tag{Ignore: true}},
		"-,noprefix":       {wantTag: Tag{Ignore: true}},
		"foo":              {wantTag: Tag{Name: "foo"}},
		"foo,inline":       {wantTag: Tag{Name: "foo", Inline: true}},
		"foo,include,":     {wantTag: Tag{Name: "foo", Include: true}},
		"foo,,noprefix":    {wantTag: Tag{Name: "foo", NoPrefix: true}},
		"FOOBAR":           {wantTag: Tag{Name: "FOOBAR"}},
		"noprefix":         {wantTag: Tag{Name: "noprefix"}},
		",noprefix":        {wantTag: Tag{NoPrefix: true}},
		",inline,noprefix": {wantTag: Tag{Inline: true, NoPrefix: true}},

		",inline,invalid": {
			wantTag: Tag{Inline: true},
			wantErr: &Error{
				TagString:   ",inline,invalid",
				Unsupported: []string{"invalid"},
			},
		},
		",extra1 ,noprefix,extra2,": {
			wantTag: Tag{NoPrefix: true},
			wantErr: &Error{
				TagString:   ",extra1 ,noprefix,extra2,",
				Unsupported: []string{"extra1 ", "extra2"},
			},
		},
	}

	for tag, tc := range tests {
		t.Run(tag, func(t *testing.T) {
			have, haveErr := ParseTag(tag)
			assert.Equal(t, tc.wantTag, have)
			if tc.wantErr == nil {
				assert.NoError(t, haveErr)
			} else {
				assert.Equal(t, tc.wantErr, errors.Unembed(haveErr))
			}
		})
	}
}

func TestParseStructField(t *testing.T) {
	type fixtureBasic struct {
		Foo string
	}
	type fixtureNamed struct {
		Foo string `env:"BAR"`
	}
	type fixtureIgnore struct {
		Foo string `env:"-" default:"some value"`
	}
	type fixtureUnexported struct {
		foo string //nolint:unused
	}

	tests := map[string]struct {
		opts  *Options
		field reflect.StructField
		want  Tag
	}{
		"basic": {
			field: reflect.TypeOf(fixtureBasic{}).Field(0),
			want:  Tag{Name: "FOO"},
		},
		"named": {
			field: reflect.TypeOf(fixtureNamed{}).Field(0),
			want:  Tag{Name: "BAR"},
		},
		"ignore": {
			field: reflect.TypeOf(fixtureIgnore{}).Field(0),
			want:  Tag{Ignore: true},
		},
		"tags only": {
			opts:  &Options{StrictTags: true},
			field: reflect.TypeOf(fixtureBasic{}).Field(0),
			want:  Tag{Ignore: true},
		},
		"unexported": {
			field: reflect.TypeOf(fixtureUnexported{}).Field(0),
			want:  Tag{Ignore: true},
		},
	}

	defaultOpts := DefaultOptions()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.opts == nil {
				tc.opts = &defaultOpts
			}

			have, _ := ParseStructField(*tc.opts, tc.field)
			assert.Equal(t, tc.want, have)
		})
	}

	t.Run("nil normalizer", func(t *testing.T) {
		have, err := ParseStructField(Options{}, reflect.TypeOf(fixtureBasic{}).Field(0))
		assert.NoError(t, err)
		assert.Equal(t, Tag{Name: "Foo"}, have)
	})

	t.Run("panic", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNormalizerEmptyName, func() {
			_, _ = ParseStructField(
				Options{
					Normalizer: NormalizerFunc(func(str string) string {
						return ""
					}),
				},
				reflect.TypeOf(fixtureBasic{}).Field(0),
			)
		})
	})
}
