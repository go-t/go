package main

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	m := make(map[interface{}]interface{})
	must(yaml.NewDecoder(os.Stdin).Decode(&m))
	must(yaml.NewEncoder(os.Stdout).Encode(m))
}
