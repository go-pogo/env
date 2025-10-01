// Copyright (c) 2025, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"reflect"
	"strings"
)

type Formatter func(name string, val any) (string, error)

// Format the name and val using a standard env format and return the resulting
// line as a string.
func Format(name string, val any) (string, error) {
	switch v := val.(type) {
	case string:
		return fmtStringValue(name, quote(v)), nil

	case Value:
		return fmtStringValue(name, quote(v.String())), nil

	case reflect.Value:
		return fmtReflectValue(name, v)

	default:
		return fmtReflectValue(name, reflect.ValueOf(v))
	}
}

// FormatShellExport the name and val using a shell compatible env format and
// return the resulting line as a string.
func FormatShellExport(name string, val any) (string, error) {
	return Format("export "+name, val)
}

func fmtStringValue(name, val string) string {
	return name + "=" + val
}

func fmtReflectValue(name string, rv reflect.Value) (string, error) {
	val, err := marshaler.Marshal(rv)
	if err != nil {
		return "", err
	}

	v := val.String()
	if rv.Kind() == reflect.String {
		v = quote(v)
	}
	return fmtStringValue(name, v), nil
}

func quote(str string) string {
	if str == "" {
		return str
	}

	isq := strings.IndexRune(str, '\'')
	idq := strings.IndexRune(str, '"')
	if isq == -1 && idq == -1 {
		return str
	}

	quot := "\""
	if isq == -1 && idq >= 0 {
		quot = "'"
	} else if isq >= 0 && idq >= 0 {
		str = strings.ReplaceAll(str, quot, "\\"+quot)
	}

	return quot + str + quot
}
