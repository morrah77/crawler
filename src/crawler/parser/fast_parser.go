package parser

import (
	"bytes"
)

type FastParser struct {
	b []byte
}

func(p *FastParser) Init(b []byte) {
	p.b = b
}

func(p *FastParser) Parse() (res [][]byte) {
	var (
		r []byte
		ok bool
	)
	for {
		if r, ok = p.ParseNext(); !ok {
			break
		}
		res = append(res, r)
	}
	return res
}

func(p *FastParser) ParseNext() (res []byte, ok bool) {
	// `<a href="https://en.wikipedia.org/wiki/UTF-8">`
	i := bytes.Index(p.b, []byte(`<a `))
	if i == -1 {
		return nil, false
	}
	j := bytes.Index(p.b[i:], []byte(`>`))
	if j == -1 {
		return nil, false
	}
	// `href="https://en.wikipedia.org/wiki/UTF-8">`
	k := bytes.Index(p.b[i:i + j], []byte(`href`))
	if k == -1 {
		return nil, false
	}
	// `="https://en.wikipedia.org/wiki/UTF-8">`
	k = bytes.Index(p.b[i + k:i + j], []byte(`=`))
	if k == -1 {
		return nil, false
	}
	// `"https://en.wikipedia.org/wiki/UTF-8">`
	l := bytes.IndexAny(p.b[i + k:i + j], ` "'`)
	if l == -1 {
		return nil, false
	}
	// `">`
	m := bytes.IndexAny(p.b[i + k + l:i + j], ` "'`)
	if m == -1 {
		return nil, false
	}
	u := p.b[i + k + l:i + m]
	u = bytes.Trim(u, ` "'`)
	if !bytes.Equal(u[:4], []byte(`http`)) {
		u = append([]byte(``), u...)
	}

	p.b = p.b[i + j:]
	return u, true
}
