package main

import (
	"fmt"

	"github.com/roeldev/go-env"
)

func main() {
	envs, n := env.Environ()

	fmt.Println("Env vars parsed:", n)
	for key, val := range envs {
		fmt.Printf("> %s: %s\n", key, val)
	}
}
