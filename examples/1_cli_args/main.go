package main

import (
	"fmt"
	"os"

	"github.com/roeldev/go-env"
)

func main() {
	args := os.Args[1:]
	envs := make(env.Map)
	_, n := env.ParseFlagArgs("e", args, envs)

	fmt.Println("Env vars parsed:", n)
	for key, val := range envs {
		fmt.Printf("> %s: %s\n", key, val)
	}
}
