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

func (gemtext Gemtext) HTML() string {
	buf := bytes.Buffer{}
	for i := 0; i < len(gemtext); i += 1 {
		switch gemtext[i].Markup {
		case Text:
			buf.Write(gemtext[i].Raw)
			buf.WriteString("\n")
		case Heading:
			buf.WriteString("<h1>")
			buf.Write(gemtext[i].Raw)
			buf.WriteString("</h1>\n")
		case SubHeading:
			buf.WriteString("<h2>")
			buf.Write(gemtext[i].Raw)
			buf.WriteString("</h2>\n")
		case SubSubHeading:
			buf.WriteString("<h3>")
			buf.Write(gemtext[i].Raw)
			buf.WriteString("</h3>\n")
		case Link:
			buf.WriteString(`<a href="`)
			parts := bytes.SplitN(gemtext[i].Raw, []byte{' '}, 3)
			_, href, text := parts[0], parts[1], parts[2]
			buf.Write(href)
			buf.WriteString(`">`)
			buf.Write(text)
			buf.WriteString("</a>\n")
		case List:
			buf.WriteString("<ul>\n")
			for ; i < len(gemtext) && gemtext[i].Markup == List; i += 1 {
				buf.WriteString("\t<li>")
				buf.Write(gemtext[i].Raw)
				buf.WriteString("</li>\n")
			}
			buf.WriteString("</ul>\n")
			i -= 1
		case Blockquote:
			buf.WriteString("<blockquote>\n")
			for ; i < len(gemtext) && gemtext[i].Markup == Blockquote; i += 1 {
				buf.Write(gemtext[i].Raw)
				buf.WriteString("\n")
			}
			buf.WriteString("</blockquote>\n")
			i -= 1
		case Preformatted:
			buf.WriteString("<pre>")
			for ; i < len(gemtext) && gemtext[i].Markup == Preformatted; i += 1 {
				buf.Write(gemtext[i].Raw)
				buf.WriteString("\n")
			}
			buf.WriteString("</pre>\n")
			i -= 1
		}
	}
	return buf.String()
}

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
