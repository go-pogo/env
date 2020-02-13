package main

import (
	"fmt"
	"strings"

	"github.com/roeldev/go-env"
)

func main() {
	reader := strings.NewReader(`
foo=bar
qux=xoo
`)

	envs := make(env.Map)
	n, err := env.Read(reader, envs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Env vars parsed:", n)
	for key, val := range envs {
		fmt.Printf("> %s: %s\n", key, val)
	}
}
