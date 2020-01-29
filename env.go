package env

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/roeldev/go-errs"
)

// Map represents a map of env variables.
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
	return Parse(os.Environ())
}

// Parse a slice of strings.
func Parse(env []string) Map {
	res := make(Map, len(env))
	for _, e := range env {
		ParseLine(e, res)
	}

	return res
}

func ParseLine(line string, dest Map) {
	// environment variables can begin with = so start splitting after the
	// first character
	split := strings.SplitAfterN(line[1:], "=", 2)
	if len(split) != 2 {
		return
	}

	key := string(line[0]) + split[0]
	dest[key[0:len(key)-1]] = split[1]
}

// Read from an `io.Reader` and parse its results.
func Read(r io.Reader) (Map, error) {
	res := make(Map, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ParseLine(scanner.Text(), res)
	}

	return res, errs.Wrap(scanner.Err())
}
