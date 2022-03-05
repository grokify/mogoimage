package main

import (
	"log"

	"github.com/grokify/goimage/primitiveutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
)

func main() {
	cfg := primitiveutil.Config{}
	_, err := flags.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Inflate()

	fmtutil.PrintJSON(cfg)

	err = cfg.Create()
	if err != nil {
		log.Fatal(err)
	}
}
