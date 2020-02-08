package env

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/roeldev/go-errs"
)

const (
	runeQuot    = 39 // '
	runeDblQuot = 34 // "
	runeHash    = 35 // #
)

// Map represents a map of env key value pairs.
type Map map[string]string

// Merge any map of strings with this Map.
func (m Map) Merge(e map[string]string) {
	for k, v := range e {
		m[k] = v
	}
}

// Environ returns a `Map` with the os' current environment variables.
func Environ() Map {
	s := os.Environ()
	m := make(Map, len(s))
	ParseSlice(s, m)
	return m
}

// Read from an `io.Reader`, parse its results and add them to the provided Map. Each line is
// cleaned before being parsed with `ParsePair`.
// It returns the number of parsed lines and any error that occurs while scanning for lines.
func Read(r io.Reader, dest Map) (int, error) {
	n := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == runeHash {
			continue // skip empty lines and comments
		}

		key, val := ParsePair(cleanLine(line))
		if key != "" && val != "" {
			dest[key] = val
			n++
		}
	}

	return n, errs.Wrap(scanner.Err())
}

// ParseSlice parses a slice of strings to a map with key value pairs. The slice should be clean,
// entries are not checked on starting/trailing whitespace or comment tags.
// It returns the number of parsed lines.
func ParseSlice(env []string, dest Map) int {
	n := 0
	for _, e := range env {
		key, val := ParsePair(e)
		dest[key] = val
		n++
	}

	return n
}

// ParsePair parses an environment variable's paired key and value and returns the separated values.
func ParsePair(pair string) (key string, val string) {
	// environment variables can begin with = so start splitting after the
	// first character
	split := strings.SplitAfterN(pair[1:], "=", 2)
	if len(split) != 2 {
		return "", ""
	}

	key = string(pair[0]) + split[0]
	// after splitting, the last character of key is '='
	// obviously we do not want to keep it
	key = key[0 : len(key)-1]

	val = split[1]
	if val == "" {
		return
	}

	last := len(val) - 1

	// remove optional double quotes from string
	if (val[0] == runeDblQuot || val[0] == runeQuot) && val[0] == val[last] {
		val = val[1:last]
	}

	return
}

func cleanLine(line string) string {
	return line
}
