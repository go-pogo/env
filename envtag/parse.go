// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"reflect"
	"strings"

	"github.com/go-pogo/errors"
)

const (
	EnvKey     = "env"
	DefaultKey = "default"
)

type Error struct {
	TagString   string
	Unsupported []string
}

func (e *Error) Error() string {
	return "error while parsing tag `" + e.TagString + "`"
}

// ParseTag parses str into a [Tag]. It will always return a usable [Tag], even
// if an error has occurred.
func ParseTag(str string) (Tag, error) {
	var t Tag
	err := parse(&t, str)
	return t, err
}

const panicNormalizerEmptyName = "envtag: Normalizer returned an empty name"

// ParseStructField uses the field's [reflect.StructTag] to look up the tag
// string according to the provided [Options]. It will always return a usable
// [Tag], even if an error has occurred.
func ParseStructField(opts Options, field reflect.StructField, prefix string) (tag Tag, err error) {
	if !field.IsExported() {
		tag.Ignore = true
		return
	}

	if str, found := field.Tag.Lookup(opts.EnvKey); found {
		// err warns for unsupported tag options, so continue parsing
		err = parse(&tag, str)
	} else if opts.StrictTags {
		tag.Ignore = true
		return
	}

	if tag.Ignore {
		return
	}
	if tag.Name == "" {
		if opts.Normalizer != nil {
			tag.Name = opts.Normalizer.Normalize(field.Name, prefix)
			if tag.Name == "" {
				panic(panicNormalizerEmptyName)
			}
		} else {
			// use field name as is, it is the callers responsibility to set a
			// valid Normalizer
			if prefix == "" {
				tag.Name = field.Name
			} else {
				tag.Name = prefix + "_" + field.Name
			}
		}
	}
	if opts.DefaultKey != "" {
		tag.Default = field.Tag.Get(opts.DefaultKey)
	}
	return
}

func parse(tag *Tag, str string) error {
	if str == "" {
		return nil
	}

	if str == "-" || strings.HasPrefix(str, "-,") {
		tag.Ignore = true
		return nil
	}

	split := strings.Split(str, ",")
	tag.Name = split[0]
	split = split[1:]

	for i, n := 0, len(split); i < n; {
		switch strings.TrimSpace(split[i]) {
		case "":
			// empty option is ignored
		case "inline":
			tag.Inline = true
		case "include":
			tag.Include = true
		default:
			// invalid options increment the index position,
			// so we end up with a slice of invalid options
			i++
			continue
		}

		// remove valid option from start of slice, all remaining options
		// shift an index position to the left. because of this, do not
		// increment i, decrease known size of slice instead
		split = append(split[:i], split[i+1:]...)
		n--
	}
	if len(split) > 0 {
		return errors.WithStack(&Error{
			TagString:   str,
			Unsupported: split,
		})
	}
	return nil
}
