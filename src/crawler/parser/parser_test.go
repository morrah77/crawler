package parser

import (
	"testing"
	"gotest.tools/assert"
	"net/url"
)

var testData = [][]byte{
	[]byte(`<html>
<head></head>
<body>
	foo<br/>
	bar\n<a href="http://www.my.com/foo/bar?baz">link</a>
	foofoofoo
	<a href="my.tmp"/>


	<a style="color:red;" href = "/my/yours/fool.html">lnk</a>
</body>
</html>`),
	[]byte(`<a href="http://www.my.com/foo/bar?baz">`),
}

var expectedResult = [][][]byte{
	{
		[]byte(`http://www.my.com/foo/bar?baz`),
		[]byte(`http://my.com/foo/my.tmp`),
		[]byte(`http://my.com/my/yours/fool.html`),
	},
	{
		[]byte(`http://www.my.com/foo/bar?baz`),
	},
}

var patterns = [][]byte{
	[]byte(`my.com`),
}

var u = url.URL{
	Scheme:`http`,
	Host:`my.com`,
	Path:`/foo`,
}

func TestNewParser(t *testing.T) {
	p, e := NewParser(testData[0], patterns, &u)
	assert.NilError(t, e, `Test NewParser(): Error should be nil!`)
	assert.Check(t, (p != nil), `Test NewParser(): Parser should not be nil!`)
}

func TestParser_Parse(t *testing.T) {
	for i, b := range testData {
		p, e := NewParser(b, patterns, &u)
		assert.NilError(t, e)
		res := p.Parse()
		assert.DeepEqual(t, res, expectedResult[i])
	}
}

func TestParser_ParseNext(t *testing.T) {
	for i, b := range testData {
		p, e := NewParser(b, patterns, &u)
		assert.NilError(t, e)
		for j, v := range expectedResult[i] {
			res, ok := p.ParseNext()
			assert.DeepEqual(t, res, v)
			assert.Equal(t, ok, j < len(expectedResult[i]))
		}
		_, ok := p.ParseNext()
		assert.Equal(t, ok, false)
	}
}
