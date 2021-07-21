package main

import (
	"bufio"
	"bytes"
	"io"
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

type Line struct {
	Raw []byte
	Markup
}

type Gemtext []Line

func ParseReader(r io.Reader) (gemtext Gemtext) {
	scanner := bufio.NewScanner(r)
	pre := false
	for scanner.Scan() {
		raw := scanner.Bytes()
		line := Line{raw, ParseMarkup(raw)}
		if line.Markup == Preformatted {
			pre = !pre
		}
		if pre {
			line.Markup = Preformatted
		}
		gemtext = append(gemtext, line)
	}
	return
}

func ParseMarkup(raw []byte) Markup {
	if bytes.HasPrefix(raw, []byte("# ")) {
		return Heading
	}
	if bytes.HasPrefix(raw, []byte("## ")) {
		return SubHeading
	}
	if bytes.HasPrefix(raw, []byte("### ")) {
		return SubSubHeading
	}
	if bytes.HasPrefix(raw, []byte("=> ")) {
		return Link
	}
	if bytes.HasPrefix(raw, []byte("* ")) {
		return List
	}
	if bytes.HasPrefix(raw, []byte("> ")) {
		return Blockquote
	}
	if bytes.HasPrefix(raw, []byte("```")) {
		return Preformatted
	}
	return Text
}
