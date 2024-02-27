// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetActiveEnvironment(t *testing.T) {
	type testCase struct {
		args       []string
		wantEnv    ActiveEnvironment
		wantRemain []string
	}
	tests := map[string]testCase{
		"nil": {},
		"empty": {
			args:       []string{""},
			wantRemain: []string{""},
		},
		"none": {
			args:       []string{"-test", flag1 + "x", flag2 + "y"},
			wantRemain: []string{"-test", flag1 + "x", flag2 + "y"},
		},
		"only flag1": {
			args:       []string{flag1, "test"},
			wantEnv:    Testing,
			wantRemain: []string{},
		},
		"only flag2": {
			args:       []string{flag2, "test"},
			wantEnv:    Testing,
			wantRemain: []string{},
		},
		"flag1 only arg": {
			args:       []string{flag1 + "=dev"},
			wantEnv:    Development,
			wantRemain: []string{},
		},
		"flag2 only arg": {
			args:       []string{flag2 + "=prod"},
			wantEnv:    Production,
			wantRemain: []string{},
		},
		"flag1 at start with remain": {
			args:       []string{flag1, "dev", "-remainder"},
			wantEnv:    Development,
			wantRemain: []string{"-remainder"},
		},
		"flag2 at start with remain": {
			args:       []string{flag2, "dev", "-remainder", flag1},
			wantEnv:    Development,
			wantRemain: []string{"-remainder", flag1},
		},
		"at end with remain": {
			args:       []string{"-remainder", flag1, "prod"},
			wantEnv:    Production,
			wantRemain: []string{"-remainder"},
		},
		"at end no value": {
			args:       []string{"-remainder", flag1},
			wantRemain: []string{"-remainder", flag1},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			argsLen := len(tc.args)
			haveEnv, haveRemain := GetActiveEnvironment(tc.args)
			assert.Equal(t, tc.wantEnv, haveEnv)
			assert.Equal(t, tc.wantRemain, haveRemain)
			assert.Len(t, tc.args, argsLen, "args length should not change")
		})
	}
}
