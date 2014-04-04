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

func init() {
	flag.Var(&color, "color", fmt.Sprintf("%s %s", "set color option", goj.Colors))
}

func main() {
	flag.Parse()

	dec := goj.NewDecoder(os.Stdin)
	dec.SetColor(color)

	f := filter()

	for {
		if err := dec.Decode(f); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println(dec)
	}

}

func filter() string {
	if len(flag.Args()) > 1 {
		return flag.Args()[0]
	}
	return ""
}
