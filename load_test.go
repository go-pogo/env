// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	restore := Environ()
	os.Clearenv()
	defer Load(restore)

	want := Value("this value is not overwritten")
	key := randKey()
	require.NoError(t, Setenv(key, want))

	assert.Nil(t, Load(Map{key: "foobar"}))
	assert.Equal(t, want, Getenv(key))
}

func TestMap_Overload(t *testing.T) {
	want := Value("foobar")
	key := randKey()
	require.NoError(t, Setenv(key, "overwrite me!"))

	assert.Nil(t, Overload(Map{key: want}))
	assert.Equal(t, want, Getenv(key))
}
