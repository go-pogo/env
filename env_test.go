package env

import (
	"os"
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
	env := Environ()
	if len(os.Environ()) != len(env) {
		t.Error(fail.Msg{
			Func: "Environ",
			Msg:  "all entries of `os.Environ()` should be present",
		})
	}
}

func TestRead(t *testing.T) {
	f, err := os.Open("test.env")
	if err != nil {
		t.Fatal(err)
	}

	have := make(Map)
	_, err = Read(f, have)
	_ = f.Close()

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
