// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

type Environment string

func (e Environment) String() string { return string(e) }

const (
	None        Environment = ""
	Development Environment = "dev"
	Testing     Environment = "test"
	Production  Environment = "prod"
)
