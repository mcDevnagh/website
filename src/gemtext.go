package main

import (
	"bufio"
	"bytes"
	"io"
)

type Attributes map[string][]byte

func (a Attributes) Write(w io.Writer) {
	for key, value := range a {
		w.Write([]byte{' '})
		io.WriteString(w, key)
		w.Write([]byte(`="`))
		w.Write(value)
		w.Write([]byte{'"'})
	}
}

type markupType uint8

const (
	blank markupType = iota
	heading
	subheading
	subsubheading
	link
	list
	blockquote
	preformatted
	text
)

var tag = [...][]byte{
	nil,
	[]byte("h1"),
	[]byte("h2"),
	[]byte("h3"),
	[]byte("a"),
	[]byte("li"),
	nil,
	nil,
	[]byte("span"),
}

var surroundTag = [...][]byte{
	nil,
	nil,
	nil,
	nil,
	nil,
	[]byte("ul"),
	[]byte("blockquote"),
	[]byte("pre"),
	nil,
}

type Markup struct {
	Raw []byte
	Attributes
	markupType
}

func (m Markup) Tag() []byte {
	return tag[m.markupType]
}

func (m Markup) SurroundTag() []byte {
	return surroundTag[m.markupType]
}

func (m Markup) Content() []byte {
	if m.markupType == link {
		return bytes.SplitN(m.Raw, []byte{' '}, 3)[2]
	}
	if m.markupType == preformatted && bytes.Equal(m.Raw, []byte("```")) {
		return nil
	}
	return m.Raw
}

func (m Markup) HTML(w io.Writer) {
	if tag := m.Tag(); tag != nil {
		w.Write([]byte{'<'})
		w.Write(tag)
		m.Attributes.Write(w)
		w.Write([]byte{'>'})
	}
	content := m.Content()
	w.Write(content)
	if tag := m.Tag(); tag != nil {
		w.Write([]byte("</"))
		w.Write(tag)
		w.Write([]byte{'>'})
	}
	if content != nil {
		w.Write([]byte{'\n'})
	}
}

type Gemtext []Markup

func (g Gemtext) HTML(w io.Writer) {
	var lastMarkup Markup
	for _, markup := range g {
		if markup.markupType != lastMarkup.markupType {
			if tag := lastMarkup.SurroundTag(); tag != nil {
				w.Write([]byte("</"))
				w.Write(tag)
				w.Write([]byte{'>'})
				w.Write([]byte{'\n'})
			}
			if tag := markup.SurroundTag(); tag != nil {
				w.Write([]byte{'<'})
				w.Write(tag)
				w.Write([]byte{'>'})
				w.Write([]byte{'\n'})
			}
		}
		markup.HTML(w)
		lastMarkup = markup
	}
}

func ParseReader(r io.Reader) (gemtext Gemtext) {
	scanner := bufio.NewScanner(r)
	pre := false
	for scanner.Scan() {
		markup := ParseLine(scanner.Bytes())
		if markup.markupType == preformatted {
			pre = !pre
		}
		if pre {
			markup.markupType = preformatted
		}
		gemtext = append(gemtext, markup)
	}
	return
}

func ParseLine(raw []byte) Markup {
	if len(raw) == 0 {
		return Markup{raw, nil, blank}
	}

	attr := make(Attributes)

	if bytes.HasPrefix(raw, []byte("# ")) {
		return Markup{raw, attr, heading}
	}
	if bytes.HasPrefix(raw, []byte("## ")) {
		return Markup{raw, attr, subheading}
	}
	if bytes.HasPrefix(raw, []byte("### ")) {
		return Markup{raw, attr, subsubheading}
	}
	if bytes.HasPrefix(raw, []byte("=> ")) {
		attr["href"] = bytes.SplitN(raw, []byte{' '}, 3)[1]
		return Markup{raw, attr, link}
	}
	if bytes.HasPrefix(raw, []byte("* ")) {
		return Markup{raw, attr, list}
	}
	if bytes.HasPrefix(raw, []byte("> ")) {
		return Markup{raw, attr, blockquote}
	}
	if bytes.HasPrefix(raw, []byte("```")) {
		return Markup{raw, attr, preformatted}
	}
	return Markup{raw, attr, text}
}
