// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"github.com/go-pogo/errors"
	"reflect"
	"strings"
)

const (
	EnvKey     = "env"
	DefaultKey = "default"

	// ignore current field + children
	ignore = "-"
	// inline means ignore name of current field, as if child fields are part
	// of parent
	inline = "inline"
	// noprefix means ignore all parent field names
	noprefix = "noprefix"
)

type Error struct {
	TagString   string
	Unsupported []string
}

func (e *Error) Error() string {
	return "error while parsing tag `" + e.TagString + "`"
}

// ParseTag parses str into a Tag. It will always return a usable Tag, even if
// an error has occurred.
func ParseTag(str string) (Tag, error) {
	var t Tag
	err := parse(&t, str)
	return t, err
}

// ParseStructField uses the field's reflect.StructTag to lookup the tag string
// according to the provided Options. It will always return a usable Tag, even
// if an error has occurred.
func ParseStructField(opts Options, field reflect.StructField) (tag Tag, err error) {
	if !field.IsExported() {
		tag.Ignore = true
		return
	}

	if str, found := field.Tag.Lookup(opts.EnvKey); found {
		// err warns for unsupported tag options, so continue parsing
		err = parse(&tag, str)
	} else if !found && opts.TagsOnly {
		tag.Ignore = true
		return
	}

	if tag.Ignore {
		return
	}
	if tag.Name == "" {
		if opts.Normalize != nil {
			tag.Name = opts.Normalize(field.Name)
		} else {
			tag.Name = field.Name
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
	if str == ignore {
		tag.Ignore = true
		return nil
	}

	split := strings.Split(str, ",")
	if split[0] == ignore {
		tag.Ignore = true
		return nil
	}

	tag.Name = split[0]
	split = split[1:]

	for i, n := 0, len(split); i < n; {
		switch strings.TrimSpace(split[i]) {
		case "":
			// empty is ignored, but is considered a valid option
		case inline:
			tag.Inline = true
		case noprefix:
			tag.NoPrefix = true
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
		return errors.WithStack(&Error{TagString: str, Unsupported: split})
	}
	return nil
}
