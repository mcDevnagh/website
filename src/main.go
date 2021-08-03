package main

import (
	"flag"
	"os"
	"website/markup"
)

var templateFile string

func init() {
	flag.StringVar(&templateFile, "t", "test.gmi",
		"Gemtext-Like file to convert to HTML, Gemtext, and/or PDF")
}

func main() {
	flag.Parse()

	f, err := os.Open(templateFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	m := markup.ParseFromGemtext(f)
	m.Gemtext(os.Stdout)
	m.HTML(os.Stdout)
}
