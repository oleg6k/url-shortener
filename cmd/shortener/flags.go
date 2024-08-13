package main

import (
	"flag"
	"net/url"
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

	runAddr, err := url.Parse(flags.RunAddr)
	if err != nil {
		panic("incorrect URL")
	}

	if runAddr.Host == "" {
		flags.RunAddr = "http://localhost" + flags.RunAddr
	}

	baseURL, err := url.Parse(flags.BaseURL)
	if err != nil {
		panic("incorrect URL")
	}

	if baseURL.Host == "" {
		flags.BaseURL = "http://localhost" + flags.BaseURL
	}
	return flags
}
