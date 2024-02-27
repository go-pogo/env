// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"strings"
)

const (
	None        ActiveEnvironment = ""
	Development ActiveEnvironment = "dev"
	Testing     ActiveEnvironment = "test"
	Production  ActiveEnvironment = "prod"

	LongFlag = "active-env"

	flag1 = "-" + LongFlag
	flag2 = "--" + LongFlag
)

// ActiveEnvironment is the active environment.
type ActiveEnvironment string

func (ae ActiveEnvironment) String() string { return string(ae) }

// GetActiveEnvironment returns the active environment and the remaining
// arguments from the provided args.
//
//	args := os.Args[1:]
//	env, args := dotenv.GetActiveEnvironment(args...)
func GetActiveEnvironment(args []string) (ActiveEnvironment, []string) {
	const flag1a = flag1 + "="
	const flag2a = flag2 + "="
	const flag1b = len(flag1a)
	const flag2b = len(flag2a)

	for i, l := 0, len(args)-1; i <= l; i++ {
		arg := args[i]
		if arg == "" {
			continue
		}
		if arg == flag1 || arg == flag2 {
			if i == l {
				// at end of args, no value set
				return "", args[:]
			}

			arg = args[i+1]
			return ActiveEnvironment(arg), append(args[:i], args[i+2:]...)
		}
		if strings.HasPrefix(arg, flag1a) {
			return ActiveEnvironment(arg[flag1b:]), append(args[:i], args[i+1:]...)
		}
		if strings.HasPrefix(arg, flag2a) {
			return ActiveEnvironment(arg[flag2b:]), append(args[:i], args[i+1:]...)
		}
	}

	return None, args
}

func GetActiveEnvironmentOr(args []string, def ActiveEnvironment) (ActiveEnvironment, []string) {
	ae, args := GetActiveEnvironment(args)
	if ae == None {
		ae = def
	}
	return ae, args
}
