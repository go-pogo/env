// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"reflect"

	"github.com/go-pogo/env/envtag"
	"github.com/go-pogo/errors"
)

// Extract environment variables names and values from the provided struct v.
func Extract(v any) (map[string]any, error) {
	return NewExtractor().Extract(v)
}

// An Extractor extracts environment variables names and values from a struct
// value.
type Extractor struct {
	TagOptions
}

// NewExtractor returns a new [Extractor].
func NewExtractor() *Extractor {
	return &Extractor{
		TagOptions: envtag.DefaultOptions(),
	}
}

// WithTagOptions sets TagOptions to the provided [TagOptions] opts.
func (ex *Extractor) WithTagOptions(opts TagOptions) *Extractor {
	ex.TagOptions = opts
	return ex
}

// Extract environment variables names and values from the provided struct v.
func (ex *Extractor) Extract(v any) (map[string]any, error) {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	if underlyingKind(rv.Type()) != reflect.Struct {
		return nil, errors.New(ErrStructExpected)
	}

	res := make(map[string]any)
	trav := &traverser{
		TagOptions:  ex.TagOptions,
		isKnownType: typeKnownByUnmarshaler,
		handleField: func(rv reflect.Value, tag envtag.Tag) (err error) {
			if rv.IsZero() && tag.Default != "" {
				if rv, err = defaultValue(rv.Type(), tag.DefaultValue()); err != nil {
					return err
				}
			}

			res[tag.Name] = rv.Interface()
			return nil
		},
	}

	if err := trav.start(rv); err != nil {
		return res, err
	}
	return res, nil
}
