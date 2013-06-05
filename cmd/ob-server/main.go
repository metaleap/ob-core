package main

import (
	"flag"

	ob "github.com/openbase/ob-core"
)

func main() {
	dirPath := flag.String("dir", ".", "usage")
	flag.Parse()
	err := ob.Init(*dirPath)
	if err != nil {
		panic(err)
	}
}
