package main

import "github.com/last-second/services/cmd"

type Foo struct {
	Key string `json:"key"`
}

func main() {
	cmd.Execute()
}
