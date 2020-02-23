package env

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/roeldev/go-fail"
)

func TestMap_Merge(t *testing.T) {
	tests := map[string]struct {
		env   Map
		merge map[string]string
		want  Map
	}{
		"append": {
			env:   Map{"foo": "bar"},
			merge: map[string]string{"qux": "xoo"},
			want:  Map{"foo": "bar", "qux": "xoo"},
		},
		"replace": {
			env:   Map{"foo": "bar", "qux": "xoo"},
			merge: map[string]string{"qux": "bar", "foo": "xoo"},
			want:  Map{"foo": "xoo", "qux": "bar"},
		},
		"merge": {
			env:   Map{"foo": "bar", "bar": "baz"},
			merge: map[string]string{"baz": "foo", "bar": "qux"},
			want:  Map{"foo": "bar", "bar": "qux", "baz": "foo"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.env.Merge(tc.merge)
			if !cmp.Equal(tc.env, tc.want) {
				t.Error(fail.Diff{
					Func: "Map.Merge",
					Have: tc.env,
					Want: tc.want,
				})
			}
		})
	}
}

func TestEnviron(t *testing.T) {
	have, n := Environ()
	if n != len(os.Environ()) || n != len(have) {
		t.Error(fail.Diff{
			Func: "Environ",
			Msg:  "all entries of `os.Environ()` should be present",
			Have: have,
			Want: os.Environ(),
		})
	}
}

func TestOpen_unexisting_file(t *testing.T) {
	have := make(Map)
	n, err := Open("doesnot.exist", have)

	if err == nil {
		t.Error(fail.Msg{
			Func: "Open",
			Msg:  "expecting an error when trying to open a file that does not exist",
		})
	}
	if n != 0 {
		t.Error(fail.Diff{
			Func: "Open",
			Msg:  "should return 0 parsed lines",
			Have: n,
			Want: 0,
		})
	}
}

func TestOpen(t *testing.T) {
	have := make(Map)
	n, err := Open("test.env", have)
	if err != nil {
		t.Fatal(fail.Err{
			Func: "Open",
			Err:  err,
		})
	}

	want := Map{
		"FOO": "bar",
		"bar": "baz",
		"qux": "#xoo",
	}

	if !cmp.Equal(have, want) {
		t.Error(fail.Diff{
			Func: "Open",
			Have: have,
			Want: want,
		})
	}
	if len(want) != n {
		t.Error(fail.Diff{
			Func: "Open",
			Msg:  "return value should be the number of parsed items",
			Have: n,
			Want: len(want),
		})
	}
}

func TestRead(t *testing.T) {
	r := strings.NewReader(`FOO=bar
bar='baz'
qux="#xoo"
`)

	have := make(Map)
	n, err := Read(r, have)
	if err != nil {
		t.Fatal(fail.Err{
			Func: "Read",
			Err:  err,
		})
	}

	want := Map{
		"FOO": "bar",
		"bar": "baz",
		"qux": "#xoo",
	}

	if !cmp.Equal(have, want) {
		t.Error(fail.Diff{
			Func: "Read",
			Have: have,
			Want: want,
		})
	}
	if len(want) != n {
		t.Error(fail.Diff{
			Func: "Read",
			Msg:  "return value should be the number of parsed items",
			Have: n,
			Want: len(want),
		})
	}
}

func TestParseFlagArgs(t *testing.T) {
	tests := map[string]struct {
		flag string
		args []string
		want Map
	}{
		"empty": {
			flag: "e",
			args: []string{},
			want: Map{},
		},
		"none": {
			flag: "e",
			args: []string{"-a", "-e", "-b=1", "-c", "2", "-e"},
			want: Map{},
		},
		"single dash": {
			flag: "e",
			args: []string{"-e=foo=bar"},
			want: Map{"foo": "bar"},
		},
		"single dash next arg": {
			flag: "env",
			args: []string{"-env", "foo=bar"},
			want: Map{"foo": "bar"},
		},
		"single double dash": {
			flag: "env",
			args: []string{"--env=qux=xoo"},
			want: Map{"qux": "xoo"},
		},
		"single double dash next arg": {
			flag: "e",
			args: []string{"--e", "qux=xoo"},
			want: Map{"qux": "xoo"},
		},
		"mixed": {
			flag: "e",
			args: []string{"-e", "foo=bar", "-e=empty", "bar", "--e=qux=xoo", "-e", "bar=baz", "--skip", "-e"},
			want: Map{"foo": "bar", "qux": "xoo", "bar": "baz"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := make(Map)
			n := ParseFlagArgs(tc.flag, tc.args, have)
			if !cmp.Equal(have, tc.want) {
				t.Error(fail.Diff{
					Func: "ParseFlagArgs",
					Have: have,
					Want: tc.want,
				})
			}
			if n != len(tc.want) {
				t.Error(fail.Diff{
					Func: "ParseFlagArgs",
					Msg:  "return value should be the number of parsed items",
					Have: n,
					Want: len(tc.want),
				})
			}
		})
	}
}

func TestParsePair(t *testing.T) {
	tests := map[string][2]string{
		"=::=::":         {"=::", "::"}, // legit windows entry
		"foo=bar":        {"foo", "bar"},
		"bar='baz'":      {"bar", "baz"},
		`qUx="xoo"`:      {"qUx", "xoo"},
		`PASSWD=$ecR3t`:  {"PASSWD", "$ecR3t"},
		"# some comment": {"", ""},
		"empty=":         {"empty", ""},
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			key, val := ParsePair(input)
			if key != want[0] || val != want[1] {
				t.Error(fail.Diff{
					Func: "ParsePair",
					Have: [2]string{key, val},
					Want: want,
				})
			}
		})
	}
}
