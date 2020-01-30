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

// Merge any map of strings with the env map.
func (m Map) Merge(e map[string]string) error {
	for k, v := range e {
		m[k] = v
	}
	return nil
}

// Environ returns a `Map` with the os' current env variables.
func Environ() Map {
	return Transform(os.Environ())
}

// Transform a slice of strings to a map with key value pairs. The slice should
// be clean, entries are not checked on starting/trailing whitespace or comment
// tags.
func Transform(env []string) Map {
	res := make(Map, len(env))
	for _, e := range env {
		key, val := SplitPair(e)
		res[key] = val
	}

	return res
}

// Read from an `io.Reader` and parse its results. Each line is cleaned before
// being parsed with `SplitPair`.
func Read(r io.Reader) (Map, error) {
	res := make(Map, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == runeHash {
			continue // skip empty lines and comments
		}

		key, val := SplitPair(cleanLine(line))
		if key != "" && val != "" {
			res[key] = val
		}
	}

	return res, errs.Wrap(scanner.Err())
}

func SplitPair(pair string) (key string, val string) {
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
