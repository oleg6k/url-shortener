package main

import (
	"flag"
)

type Flags struct {
	RunAddr string
	BaseURL string
}

func ParseFlags() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flags.BaseURL, "b", ":8080", "base address and port to shortener results")
	flag.Parse()

	return flags
}
