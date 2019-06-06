package crawler

import (
	"testing"
	"bytes"
	"gotest.tools/assert"
	"net/http/httptest"
	"net/http"
)

var testData = map[string][]byte{
	`/`: []byte(`<!DOCTYPE html>
<html>
<head>
</head>
<body>
<a href="/doc">Documents</a>
<a href="/pkg">Packages</a>
</body>`),
	`/doc`: []byte(`<!DOCTYPE html>
<html>
<head>
</head>
<body>
<a href="/">Main</a>
<a href="/pkg">Packages</a>
</body>`),
	`/pkg`: []byte(`<!DOCTYPE html>
<html>
<head>
</head>
<body>
<a href="/">Main</a>
<a href="/doc">Documents</a>
<a href="/pkg">Packages</a>
</body>`),
}

func expectedResult(prefix string) string {
	return prefix + `/doc
` + prefix + `/
` + prefix + `/pkg
`
}

func TestNewCrawler(t *testing.T) {
	w := bytes.NewBuffer(make([]byte, 100))
	c := NewCrawler(w)
	assert.Equal(t, c.output, w)
}

func TestCrawler_Crawl(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.RequestURI {
		case `/`:
			writer.Write(testData[`/`])
		case `/doc`:
			writer.Write(testData[`/doc`])
		case `/pkg`:
			writer.Write(testData[`/pkg`])
		default:
			writer.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	u := srv.URL + `/doc`
	w := bytes.NewBuffer(make([]byte, 100))
	c := NewCrawler(w)
	c.Run(u)
	assert.Equal(t, w.String(), expectedResult(srv.URL))
}
