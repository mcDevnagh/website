package main

import (
	"flag"
	"os"
	"website/markup"
)

var templateFile string
var html bool
var gemtext bool

func init() {
	flag.StringVar(&templateFile, "t", "",
		"Gemtext-Like file to convert to HTML, Gemtext, and/or PDF. Exclude to read from STDIN")

	flag.BoolVar(&html, "h", false, "Convert to HTML and output to STDOUT")
	flag.BoolVar(&gemtext, "g", false, "Convert to Gemtext and output to STDOUT")
}

func main() {
	flag.Parse()

	var m markup.Markups
	if templateFile == "" {
		m = markup.ParseFromGemtext(os.Stdin)
	} else {
		f, err := os.Open(templateFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		m = markup.ParseFromGemtext(f)
	}

	if html {
		m.HTML(os.Stdout)
	}
	if gemtext {
		m.Gemtext(os.Stdout)
	}
}
