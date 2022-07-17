// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"strings"
	"unicode"

	"github.com/go-pogo/errors"
)

// Errors returned by Parse.
const (
	ParseError         errors.Kind = "parse error"
	ErrInvalidFormat   errors.Msg  = "invalid format"
	ErrMissingEndQuote errors.Msg  = "missing end quote"
	ErrEmptyKey        errors.Msg  = "empty key"

	commentHash = '#'
	quoteSingle = '\''
	quoteDouble = '"'
)

// Parse parses a string containing a possible key value pair. Any whitespace
// at the start and/or end of str is trimmed.
// It returns an empty key value pair when the provided str, after trimming,
// begins with #.
func Parse(str string) (string, Value, error) {
	str = strings.TrimSpace(str)
	if str == "" || str[0] == commentHash {
		return "", "", nil
	}
	return parse(str, true)
}

func parseAndStore(dest map[string]Value, str string, stripExport bool) (string, Value, error) {
	k, v, err := parse(str, stripExport)
	if err == nil {
		dest[k] = v
	}
	return k, v, err
}

func parse(str string, stripExport bool) (string, Value, error) {
	parts := strings.SplitAfterN(str, "=", 2)
	if len(parts) != 2 {
		return "", "", errors.WithKind(ErrInvalidFormat, ParseError)
	}

	n := len(parts[0]) - 1
	if n == 0 {
		return "", "", errors.WithKind(ErrEmptyKey, ParseError)
	}

	// strip `=` from end of key part
	key := parts[0][:n]

	if stripExport &&
		n > 6 &&
		strings.HasPrefix(key, "export") &&
		unicode.IsSpace(rune(key[6])) {
		// strip `export ` from start and optional whitespace from end
		key = strings.TrimSpace(key[6:])
	} else if unicode.IsSpace(rune(key[n-1])) {
		// strip any possible extra whitespace from end
		key = strings.TrimSpace(key[:n-1])
	}
	if key == "" {
		return "", "", errors.WithKind(ErrEmptyKey, ParseError)
	}

	val := parts[1]
	if val == "" {
		return key, Value(val), nil
	}

	if unicode.IsSpace(rune(val[0])) {
		val = strings.TrimSpace(val[1:])
	}

	var err error
	switch true {
	case val[0] == quoteSingle:
		val, err = parseQuotedValue(val[1:], quoteSingle)
		if err != nil {
			return "", "", errors.WithKind(err, ParseError)
		}

	case val[0] == quoteDouble:
		val, err = parseQuotedValue(val[1:], quoteDouble)
		if err != nil {
			return "", "", errors.WithKind(err, ParseError)
		}

	default:
		i := strings.IndexRune(val, commentHash)
		if i == 0 {
			val = ""
		} else if i > 0 {
			val = val[:i-1]
		}
	}

	return key, Value(val), nil
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
				// quote is not escaped, we've reached to end of the value
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
