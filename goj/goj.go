package main

import (
	"flag"
	"fmt"
	"github.com/bhenderson/goj"
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
	if debug {
		log.Println(filter)
	}

	if err != nil {
		log.Fatalln(err)
	}

	dec := goj.NewDecoder(files...)
	dec.SetColor(color)

	out := dec.Decode(filter)
	retval := 0

	for val := range out {
		if val.Error != nil {
			fmt.Println(val.Error)
		}
		if diff {
			b, _ := goj.Diff(val, <-out)
			if len(b) > 0 {
				retval = 1
			}
			fmt.Printf("%s", b)
		} else {
			fmt.Println(val)
		}
	}

	os.Exit(retval)
}
