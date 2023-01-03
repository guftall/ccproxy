package main

import (
	"flag"

	"github.com/guftall/ccproxy/pkg/server"
)

func main() {
	flag.Parse()

	s := &server.Server{}
	s.Start()
}
