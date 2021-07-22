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

	fmt.Print(ParseReader(f).HTML())
}
