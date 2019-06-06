package parser

import (
	"regexp"
	"strings"
	"bytes"
	"net/url"
)

type Parser struct {
	body []byte
	parsed bool
	matches map[string] struct{}
	patterns []*regexp.Regexp
	url *url.URL
}

var rawPatterns = [][]byte{
	[]byte(`\<a[^<>]*href\s*=\s*\"?([^'"<>]*)\"?[^>]*\>`),
}

func NewParser(body []byte, domains [][]byte, url *url.URL) (p *Parser, err error) {
	p = &Parser{
		body:body,
		matches:make(map[string]struct{}),
		patterns:make([]*regexp.Regexp, 0),
		url:url,
	}
	for _, pt := range rawPatterns {
		re, err := regexp.Compile(string(pt))
		if err != nil {
			return nil, err
		}
		p.patterns = append(p.patterns, re)
	}
	for _, pt := range domains {
		s := strings.Trim(string(pt), `/ `) + `/\S*`
		re, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		p.patterns = append(p.patterns, re)
	}
	return p, err
}

func(p *Parser) Init(b []byte) {
	p.parsed = false
	p.matches = make(map[string]struct{})
	p.body = b
}

// Parse() returns all references to appropriate domain pages
//pattern := `<a href="`
//p1 := `<a href="#Identifiers">`
//p2 := `<a href="https://en.wikipedia.org/wiki/UTF-8">`
//p3 := domain
func(p *Parser) Parse() (matches [][]byte) {
	for _, re := range p.patterns {
		found := re.FindAllSubmatch(p.body, 1000)
		for _, f := range found {
			if len(f) > 1 {
				for i := 1; i < len(f); i++ {
					if bytes.HasPrefix(f[i], []byte(`http`)) {
						if (bytes.Index(f[i], []byte(p.url.Host)) > 6) {
							matches = append(matches, f[i])
						}
					}  else if f[i][0] == []byte(`/`)[0] {
						f[i] = append([]byte(p.url.Scheme + `://` + p.url.Host), f[i]...)
						matches = append(matches, f[i])
					} else {
						f[i] = append([]byte(p.url.Scheme + `://` + p.url.Host + p.url.EscapedPath() + `/`), f[i]...)
						matches = append(matches, f[i])
					}
				}

			}
		}
	}
	p.parsed = true
	return matches
}

func(p *Parser) ParseNext() (b []byte, ok bool) {
	if ! p.parsed {
		for _, ms := range p.Parse() {
			if _, ok := p.matches[string(ms)]; !ok {
				p.matches[string(ms)] = struct{}{}
			}
		}
	}
	for m, _ := range p.matches {
		delete(p.matches, m)
		return []byte(m), true
	}
	return nil, false
}
