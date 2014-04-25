package goj

import (
	"flag"
	"os"
)

func ParseFlags() (filter string, files []File, err error) {
	flag.Parse()
	i := 0
	for ; i < flag.NArg(); i++ {
		p := flag.Arg(i)
		if p == "--" {
			// skip
		} else if p == "-" {
			files = append(files, os.Stdin)
		} else {
			if f, err := os.Open(p); err != nil {
				if os.IsNotExist(err) {
					// reset error
					break
				} else {
					// some other error
					return filter, files, err
				}
			} else {
				files = append(files, f)
			}
		}
	}
	//
	if len(files) == 0 {
		files = append(files, os.Stdin)
	}
	// filter is the first non file argument.
	filter = flag.Arg(i)
	return
}
