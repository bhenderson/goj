package main

import (
	"fmt"
	"github.com/bhenderson/goj"
	"io"
	"log"
	"os"
)

func main() {
	dec := goj.NewDecoder(os.Stdin)

	for {
		if err := dec.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println(dec)
	}

}
