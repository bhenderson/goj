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
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	return ""
}
