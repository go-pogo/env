// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testEnviron() Map {
	restore := Environ()
	os.Clearenv()
	return restore
}

func restoreEnviron(restore Map) {
	os.Clearenv()
	if err := Load(restore); err != nil {
		panic(err)
	}
}

func TestLoad(t *testing.T) {
	restore := testEnviron()
	defer restoreEnviron(restore)

	want := Value("this value is not overwritten")
	key := randKey()
	require.NoError(t, Setenv(key, want))

	assert.Nil(t, Load(Map{key: "foobar"}))
	assert.Equal(t, want, Getenv(key))
}

func TestMap_Overload(t *testing.T) {
	restore := testEnviron()
	defer restoreEnviron(restore)

	want := Value("foobar")
	key := randKey()
	require.NoError(t, Setenv(key, "overwrite me!"))

	assert.Nil(t, Overload(Map{key: want}))
	assert.Equal(t, want, Getenv(key))
}

func randKey() string {
	return "somewhat_random_key_" + strconv.FormatInt(time.Now().Unix(), 10)
}
