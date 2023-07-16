// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"bufio"
	"bytes"
	"io"
	"unicode"

	"github.com/go-pogo/errors"
)

type Scanner interface {
	Err() error
	Bytes() []byte
	Text() string
	Scan() bool
	KeyValue() (string, Value, error)
}

var _ Scanner = new(scanner)

type scanner struct {
	*bufio.Scanner
}

// NewScanner returns a new Scanner which wraps a bufio.Scanner that reads from
// io.Reader r. The split function defaults to ScanLines.
// Successive calls to the Scan method will (just like bufio.Scanner) step
// through the 'tokens' of the read bytes, skipping the bytes between the
// tokens.
func NewScanner(r io.Reader) Scanner {
	if r == nil {
		return new(nilScanner)
	}

	s := bufio.NewScanner(r)
	s.Split(ScanLines)
	return &scanner{s}
}

// ScanLines is a bufio.SplitFunc that returns each line of text using
// bufio.ScanLines. Additionally, any leading or trailing whitespace is
// stripped from the token result. Lines that start with a #, after all leading
// whitespace is stripped, are treated as comments and result in an empty token
// result.
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// bufio.ScanLines always returns a nil error
	advance, token, _ = bufio.ScanLines(data, atEOF)
	if len(token) == 0 {
		return advance, token, err
	}
	if unicode.IsSpace(rune(token[0])) {
		token = bytes.TrimLeftFunc(token, unicode.IsSpace)
		if len(token) == 0 {
			return advance, token, nil
		}
	}
	if token[0] == '#' {
		return advance, token[:0], nil
	}

	token = bytes.TrimRightFunc(token, unicode.IsSpace)
	return advance, token, nil
}

func (s *scanner) Err() error { return errors.WithStack(s.Scanner.Err()) }

// Scan advances the scanner to the next token, which will then be available
// through the Bytes or Text method. The token is guaranteed to not be empty.
// See bufio.Scanner for additional details.
func (s *scanner) Scan() bool {
	ok := s.Scanner.Scan()
	// lookup next token when an empty line is encountered
	for ok && len(s.Scanner.Bytes()) == 0 {
		ok = s.Scanner.Scan()
	}
	return ok
}

func (s *scanner) KeyValue() (string, Value, error) {
	line := s.Scanner.Text()
	if line == "" {
		return "", "", nil
	}

	return parse(line)
}

type nilScanner struct{}

func (s *nilScanner) Err() error                       { return nil }
func (s *nilScanner) Bytes() []byte                    { return nil }
func (s *nilScanner) Text() string                     { return "" }
func (s *nilScanner) Scan() bool                       { return false }
func (s *nilScanner) KeyValue() (string, Value, error) { return "", "", nil }
