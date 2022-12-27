// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"net/url"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func ExampleUnmarshal() {
	type Envs struct {
		Foo string
		Bar struct {
			Url url.URL
		}
		Timeout time.Duration `default:"10s"`
		//Ip      net.IP
	}

	var data = `
FOO=bar
# ignore me
BAR_URL=http://example.com
IP=192.168.1.1`

	var envs Envs
	if err := Unmarshal([]byte(data), &envs); err != nil {
		panic(err)
	}

	spew.Dump(envs)
	// Output:
	// (env.Envs) {
	//  Foo: (string) (len=3) "bar",
	//  Bar: (struct { Url url.URL }) {
	//   Url: (url.URL) http://example.com
	//  },
	//  Timeout: (time.Duration) 10s
	// }
}

func TestUnmarshal(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type fixture struct {
			Foo    string
			Ignore bool `env:"-"`
		}

		var have fixture
		haveErr := Unmarshal([]byte("FOO=bar\nIGNORE=true"), &have)
		assert.Exactly(t, fixture{Foo: "bar"}, have)
		assert.Nil(t, haveErr)
	})
}
