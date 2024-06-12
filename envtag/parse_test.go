// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"reflect"
	"testing"

	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	tests := map[string]struct {
		wantTag Tag
		wantErr error
	}{
		"":             {wantTag: Tag{}},
		"-":            {wantTag: Tag{Ignore: true}},
		"-,inline":     {wantTag: Tag{Ignore: true}},
		"foo":          {wantTag: Tag{Name: "foo"}},
		"foo,inline":   {wantTag: Tag{Name: "foo", Inline: true}},
		"foo,include,": {wantTag: Tag{Name: "foo", Include: true}},
		"FOOBAR":       {wantTag: Tag{Name: "FOOBAR"}},
		",inline":      {wantTag: Tag{Inline: true}},

		",inline,invalid": {
			wantTag: Tag{Inline: true},
			wantErr: &Error{
				TagString:   ",inline,invalid",
				Unsupported: []string{"invalid"},
			},
		},
		",extra1 ,inline,extra2,": {
			wantTag: Tag{Inline: true},
			wantErr: &Error{
				TagString:   ",extra1 ,inline,extra2,",
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
		opts   *Options
		field  reflect.StructField
		prefix string
		want   Tag
	}{
		"field name": {
			field: reflect.TypeOf(fixtureBasic{}).Field(0),
			want:  Tag{Name: "FOO"},
		},
		"field name with prefix": {
			field:  reflect.TypeOf(fixtureBasic{}).Field(0),
			prefix: "PREFIX",
			want:   Tag{Name: "PREFIX_FOO"},
		},
		"field name with prefix and without normalizer": {
			opts:   &Options{Normalizer: nil},
			field:  reflect.TypeOf(fixtureBasic{}).Field(0),
			prefix: "PREFIX",
			want:   Tag{Name: "PREFIX_Foo"},
		},
		"tag name": {
			field: reflect.TypeOf(fixtureNamed{}).Field(0),
			want:  Tag{Name: "BAR"},
		},
		"tag name with prefix": {
			field:  reflect.TypeOf(fixtureNamed{}).Field(0),
			prefix: "PREFIX",
			want:   Tag{Name: "BAR"},
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

			have, _ := ParseStructField(*tc.opts, tc.field, tc.prefix)
			assert.Equal(t, tc.want, have)
		})
	}

	t.Run("nil normalizer", func(t *testing.T) {
		have, err := ParseStructField(Options{}, reflect.TypeOf(fixtureBasic{}).Field(0), "")
		assert.NoError(t, err)
		assert.Equal(t, Tag{Name: "Foo"}, have)
	})

	t.Run("panic", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNormalizerEmptyName, func() {
			_, _ = ParseStructField(
				Options{
					Normalizer: NormalizerFunc(func(_, _ string) string {
						return ""
					}),
				},
				reflect.TypeOf(fixtureBasic{}).Field(0),
				"",
			)
		})
	})
}
