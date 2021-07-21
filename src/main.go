package main

import (
	"flag"
	"fmt"
	"os"
)

var templateFile string

func init() {
	flag.StringVar(&templateFile, "t", "",
		"Gemtext-Like file to convert to HTML, Gemtext, and/or PDF")
}

func main() {
	flag.Parse()

	f, err := os.Open(templateFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, gt := range ParseReader(f) {
		fmt.Printf("%s\t%s\n", gt.Markup, gt.Raw)
	}
}
