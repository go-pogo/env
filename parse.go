// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/rawconv"
)

const (
	ErrInvalidFormat   errors.Msg = "invalid format"
	ErrMissingEndQuote errors.Msg = "missing end quote"
	ErrEmptyKey        errors.Msg = "empty key"
)

// Value is an alias of rawconv.Value.
type Value = rawconv.Value

type NamedValue struct {
	Name  string
	Value Value
}

func (nv NamedValue) GoString() string {
	return `env.NamedValue(` + nv.Name + `="` + nv.Value.String() + `")`
}

// Parse parses a string containing a possible key value pair. Any whitespace
// at the start and/or end of str is trimmed. It returns an empty NamedValue
// when the provided str, after trimming, begins with #.
func Parse(str string) (NamedValue, error) {
	str = strings.TrimSpace(str)
	if str == "" || str[0] == '#' {
		return NamedValue{}, nil
	}
	return parse(str)
}

type ParseError struct {
	Err error
	Str string
}

func (e *ParseError) Unwrap() error { return e.Err }

func (e *ParseError) Error() string {
	return fmt.Sprintf("error while parsing `%s`", e.Str)
}

func parse(str string) (NamedValue, error) {
	parts := strings.SplitAfterN(str, "=", 2)
	if len(parts) != 2 {
		return NamedValue{}, errors.WithStack(&ParseError{
			Err: ErrInvalidFormat,
			Str: str,
		})
	}

	n := len(parts[0]) - 1
	if n == 0 {
		return NamedValue{}, errors.WithStack(&ParseError{
			Err: ErrEmptyKey,
			Str: str,
		})
	}

	// strip `=` from end of key part
	key := parts[0][:n]
	if n > 6 &&
		strings.HasPrefix(key, "export") &&
		unicode.IsSpace(rune(key[6])) {
		// strip `export ` from start and optional whitespace from end
		key = strings.TrimSpace(key[6:])
	} else if unicode.IsSpace(rune(key[n-1])) {
		// strip any possible extra whitespace from end
		key = strings.TrimSpace(key[:n-1])
	}
	if key == "" {
		return NamedValue{}, errors.WithStack(&ParseError{
			Err: ErrEmptyKey,
			Str: str,
		})
	}

	val := parts[1]
	if val == "" {
		return NamedValue{key, Value(val)}, nil
	}

	if unicode.IsSpace(rune(val[0])) {
		val = strings.TrimSpace(val[1:])
	}

	var err error
	switch true {
	case val[0] == '\'':
		val, err = parseQuotedValue(val[1:], '\'')
		if err != nil {
			return NamedValue{}, errors.WithStack(&ParseError{
				Err: err,
				Str: str,
			})
		}

	case val[0] == '"':
		val, err = parseQuotedValue(val[1:], '"')
		if err != nil {
			return NamedValue{}, errors.WithStack(&ParseError{
				Err: err,
				Str: str,
			})
		}

	default:
		i := strings.IndexRune(val, '#')
		if i == 0 {
			val = ""
		} else if i > 0 {
			val = val[:i-1]
		}
	}

	return NamedValue{key, Value(val)}, nil
}

func parseQuotedValue(val string, q rune) (string, error) {
	// first rune q is already stripped from val
	if len(val) == 1 && rune(val[0]) == q {
		return "", nil
	}
	if !strings.ContainsRune(val, q) {
		return "", ErrMissingEndQuote
	}

	var e bool
	var s strings.Builder
	s.Grow(len(val))

loop:
	for _, char := range val {
		switch char {
		case q:
			if !e {
				// quote is not escaped, we've reached the end of the value
				break loop
			}

			e = false
			s.WriteRune(char)

		case '\\':
			if e {
				// slash is escaped which means we'll keep it
				s.WriteRune(char)
			}
			e = !e

		default:
			e = false
			s.WriteRune(char)
		}
	}

	return s.String(), nil
}
