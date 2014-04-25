package main

import (
	"flag"
	"fmt"
	"github.com/bhenderson/goj"
	"io"
	"log"
	"os"
)

var (
	color = goj.ColorAuto
	debug bool
	diff  bool
)

func init() {
	flag.Var(&color, "color", fmt.Sprintf("%s %s", "set color option", goj.Colors))
	flag.BoolVar(&debug, "debug", false, "set debugging")
	flag.BoolVar(&diff, "diff", false, "set diff option")
}

func main() {
	filter, files, err := goj.ParseFlags()

	if err != nil {
		log.Fatalln(err)
	}

	dec := goj.NewDecoder(files...)
	dec.SetColor(color)

	for val := range dec.Decode(filter, diff) {
		fmt.Println(val)
	}
}
