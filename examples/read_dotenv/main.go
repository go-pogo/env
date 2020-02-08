package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/roeldev/go-env"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	path, err := filepath.Abs("test.env")
	if err != nil {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}

	defer f.Close()

	envs := make(env.Map)
	n, err := env.Read(f, envs)
	if err != nil {
		return
	}

	fmt.Println("Env vars parsed:", n)
	for key, val := range envs {
		fmt.Printf("> %s: %s\n", key, val)
	}
}
