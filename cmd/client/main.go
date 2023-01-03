package main

import (
	"flag"

	"github.com/guftall/ccproxy/pkg/client"
)

func main() {
	flag.Parse()

	c := client.Client{}
	c.Start()
}
