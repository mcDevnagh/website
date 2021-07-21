package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Markup byte

const (
	Text Markup = iota
	Heading
	SubHeading
	SubSubHeading
	Link
	List
	Blockquote
	Preformatted
)

var markup = [...]string{
	"Text",
	"Heading",
	"SubHeading",
	"SubSubHeading",
	"Link",
	"List",
	"Blockquote",
	"Preformatted",
}

func (m Markup) String() string {
	return markup[m]
}

type Gemtext struct {
	Line string
	Markup
}

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
		fmt.Printf("%s\t%s\n", gt.Markup, gt.Line)
	}
}

func ParseReader(r io.Reader) (gemtext []Gemtext) {
	scanner := bufio.NewScanner(r)
	pre := false
	for scanner.Scan() {
		line := scanner.Text()
		gt := Gemtext{line, ParseLine(line)}
		if gt.Markup == Preformatted {
			pre = !pre
		}
		if pre {
			gt.Markup = Preformatted
		}
		gemtext = append(gemtext, gt)
	}
	return
}

func ParseLine(line string) Markup {
	if strings.HasPrefix(line, "# ") {
		return Heading
	}
	if strings.HasPrefix(line, "## ") {
		return SubHeading
	}
	if strings.HasPrefix(line, "### ") {
		return SubSubHeading
	}
	if strings.HasPrefix(line, "=> ") {
		return Link
	}
	if strings.HasPrefix(line, "* ") {
		return List
	}
	if strings.HasPrefix(line, "> ") {
		return Blockquote
	}
	if strings.HasPrefix(line, "```") {
		return Preformatted
	}
	return Text
}
