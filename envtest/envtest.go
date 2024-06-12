// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtest

import (
	"os"

	"github.com/go-pogo/env"
)

type Snapshot struct {
	envs env.Map
}

func Prepare(testenvs env.Map) *Snapshot {
	r := Snapshot{envs: env.Environ()}
	os.Clearenv()

	if len(testenvs) > 0 {
		if err := env.Load(testenvs); err != nil {
			panic(err)
		}
	}

	return &r
}

func (r *Snapshot) Restore() {
	os.Clearenv()
	if err := env.Load(r.envs); err != nil {
		panic(err)
	}
}
