package env

import (
	"fmt"
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
	m, n := Environ()
	want := os.Environ()

	if n != len(want) || n != len(m) {
		have := make([]string, 0, n)

		// loop trough the original environ list as this is in the
		// "proper" order
		for _, w := range want {
			k, _ := ParsePair(w)
			v, ok := m[k]
			if !ok {
				v = "!!MISSING!!"
			}

			have = append(have, k+"="+v)
		}

		t.Error(fail.Diff{
			Func: "Environ",
			Msg:  fmt.Sprintf("all entries of `os.Environ()` should be present [%d/%d]", n, len(want)),
			Have: have,
			Want: want,
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
		flag    string
		input   []string
		wantRes []string
		wantMap Map
	}{
		"empty": {
			flag:    "e",
			input:   []string{},
			wantRes: []string{},
			wantMap: Map{},
		},
		"none": {
			flag:    "e",
			input:   []string{"-a", "-e", "-b=1", "-c", "2", "-e"},
			wantRes: []string{"-a", "-b=1", "-c", "2"},
			wantMap: Map{},
		},
		"single dash": {
			flag:    "e",
			input:   []string{"-e=foo=bar"},
			wantRes: []string{},
			wantMap: Map{"foo": "bar"},
		},
		"single dash next arg": {
			flag:    "env",
			input:   []string{"-env", "foo=bar"},
			wantRes: []string{},
			wantMap: Map{"foo": "bar"},
		},
		"single double dash": {
			flag:    "env",
			input:   []string{"--env=qux=xoo"},
			wantRes: []string{},
			wantMap: Map{"qux": "xoo"},
		},
		"single double dash next arg": {
			flag:    "e",
			input:   []string{"--e", "qux=xoo"},
			wantRes: []string{},
			wantMap: Map{"qux": "xoo"},
		},
		"mixed": {
			flag:    "e",
			input:   []string{"-e", "foo=bar", "-e=empty", "bar", "--e=qux=xoo", "-e", "bar=baz", "--skip", "-e"},
			wantRes: []string{"-e=empty", "bar", "--skip"},
			wantMap: Map{"foo": "bar", "qux": "xoo", "bar": "baz"},
		},
		"multi mixed": {
			flag:    "e",
			input:   []string{"-e", "foo=bar", "empty=", "bar", "--e=qux=xoo", "nop=nop", "-e", "bar=baz", "--skip", "-e"},
			wantRes: []string{"bar", "nop=nop", "--skip"},
			wantMap: Map{"foo": "bar", "empty": "", "qux": "xoo", "bar": "baz"},
		},
		"lookahead": {
			flag:    "e",
			input:   []string{"-e", "-t", "foo=bar", "-e", "qux=xoo", "empty=", "baz"},
			wantRes: []string{"-t", "foo=bar", "baz"},
			wantMap: Map{"qux": "xoo", "empty": ""},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			haveMap := make(Map)
			haveRes, n := ParseFlagArgs(tc.flag, tc.input, haveMap)
			if !cmp.Equal(haveMap, tc.wantMap) {
				t.Error(fail.Diff{
					Func: "ParseFlagArgs",
					Msg:  "destination map should include all parsed env. vars",
					Have: haveMap,
					Want: tc.wantMap,
				})
			}
			if !cmp.Equal(haveRes, tc.wantRes) {
				t.Error(fail.Diff{
					Func: "ParseFlagArgs",
					Msg:  "return value should be without parsed env. vars",
					Have: strings.Join(haveRes, " "),
					Want: strings.Join(tc.wantRes, " "),
				})
			}
			if n != len(tc.wantMap) {
				t.Error(fail.Diff{
					Func: "ParseFlagArgs",
					Msg:  "second return value should be the number of parsed items",
					Have: n,
					Want: len(tc.wantMap),
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
