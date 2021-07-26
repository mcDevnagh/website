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
	[]byte("p"),
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

var markupLen = [...]int{0, 2, 3, 4, 3, 2, 2, 0, 0}

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
	if m.markupType == preformatted && bytes.Equal(m.Raw, []byte("```")) {
		return nil
	}

	space := []byte{' '}

	var content []byte
	attrCount := len(m.Attributes)
	if m.markupType == link {
		content = bytes.SplitN(m.Raw, space, 3)[2]
		attrCount -= 1
	} else {
		content := m.Raw[markupLen[m.markupType]:]
	}
	for i := 0; i < attrCount; i++ {
		content = bytes.TrimSpace(content)
		content = content[:bytes.LastIndex(content, space)]
	}
	return content
}

func writeTag(w io.Writer, tag []byte, attr Attributes, closeTag bool) {
	if tag != nil {
		w.Write([]byte{'<'})
		if closeTag {
			w.Write([]byte{'/'})
		}
		w.Write(tag)
		attr.Write(w)
		w.Write([]byte{'>'})
	}
}

func (m Markup) HTML(w io.Writer) {
	writeTag(w, m.Tag(), m.Attributes, false)
	content := m.Content()
	w.Write(content)
	writeTag(w, m.Tag(), nil, true)
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
				writeTag(w, tag, nil, true)
				w.Write([]byte{'\n'})
			}
			if tag := markup.SurroundTag(); tag != nil {
				writeTag(w, tag, nil, false)
				w.Write([]byte{'\n'})
			}
		}
		markup.HTML(w)
		lastMarkup = markup
	}
	if tag := lastMarkup.SurroundTag(); tag != nil {
		writeTag(w, lastMarkup.SurroundTag(), nil, true)
		w.Write([]byte{'\n'})
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

var selectors = []byte{'#', '.'}
var attrKeys = map[byte]string{
	'#': "id",
	'.': "class",
}

func ParseLine(raw []byte) Markup {
	if len(raw) == 0 {
		return Markup{raw, nil, blank}
	}
	if bytes.HasPrefix(raw, []byte("```")) {
		return Markup{raw, nil, preformatted}
	}

	var attr Attributes
	for words := bytes.Split(bytes.TrimSpace(raw), []byte{' '}); len(words) > 0; words = words[:len(words)-1] {
		word := words[len(words)-1]
		if len(word) <= 1 || !bytes.Contains(selectors, word[0:1]) {
			break
		}
		if attr == nil {
			attr = make(map[string][]byte)
		}
		attr[attrKeys[word[0]]] = word[1:]
	}

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
		if attr == nil {
			attr = make(Attributes)
		}
		attr["href"] = bytes.SplitN(raw, []byte{' '}, 3)[1]
		return Markup{raw, attr, link}
	}
	if bytes.HasPrefix(raw, []byte("* ")) {
		return Markup{raw, attr, list}
	}
	if bytes.HasPrefix(raw, []byte("> ")) {
		return Markup{raw, attr, blockquote}
	}
	return Markup{raw, attr, text}
}
