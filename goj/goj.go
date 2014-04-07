package main

import (
	"flag"
	"fmt"
	"github.com/bhenderson/goj"
	"io"
	"log"
	"os"
)

var color = goj.ColorAuto
var filter string
var debug bool

func init() {
	flag.Var(&color, "color", fmt.Sprintf("%s %s", "set color option", goj.Colors))
	flag.BoolVar(&debug, "debug", false, "set debugging")
}

func main() {
	parseFlags()

	var prev *goj.Decoder
	var diff bool
	dec := goj.NewDecoder(os.Stdin)
	dec.SetColor(color)

	for {
		if err := dec.Decode(filter); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if prev != nil {
			b, _ := goj.Diff(prev, dec)
			fmt.Println(string(b))
			diff = true
		}
		prev = dec.Copy()
	}

	if !diff {
		fmt.Println(dec)
	}
}

func parseFlags() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		filter = flag.Args()[0]
		if debug {
			log.Println(filter)
		}
	}
}
