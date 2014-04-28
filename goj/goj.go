package main

import (
	"flag"
	"fmt"
	"github.com/bhenderson/goj"
	"log"
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
	if debug {
		log.Println(filter)
	}

	if err != nil {
		log.Fatalln(err)
	}

	dec := goj.NewDecoder(files...)
	dec.SetColor(color)

	out := dec.Decode(filter)

	for val := range out {
		if diff {
			b, _ := goj.Diff(val, <-out)
			fmt.Println(string(b))
		} else {
			fmt.Println(val)
		}
	}
}
