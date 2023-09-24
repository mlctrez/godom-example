package main

import (
	"github.com/mlctrez/godom-example/example"
	"github.com/mlctrez/godom/app/ctx"
)

func main() {
	ctx.Run(example.New())
}
