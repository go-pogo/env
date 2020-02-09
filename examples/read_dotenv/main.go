package main

import (
	"fmt"
	"path/filepath"

	"github.com/roeldev/go-env"
)

func main() {
	path, err := filepath.Abs("test.env")
	if err != nil {
		panic(err)
	}

	envs := make(env.Map)
	n, err := env.Open(path, envs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Env vars parsed:", n)
	for key, val := range envs {
		fmt.Printf("> %s: %s\n", key, val)
	}
}
