package crawler

import (
	"io"
	"net/http"
	"io/ioutil"
	"net/url"
	"errors"
	"time"
	"sync"
	"crawler/parser"
	"crawler/common"
)

type uriRequest struct {
	uri string
	respChannel chan bool
}

type Crawler struct {
	output            io.Writer
	setOutputChannel chan []byte
	visitedUris       map[string]struct{}
	setVisitedChannel chan string
	getVisitedChannel chan uriRequest
	urisToVisit       chan string
	stop              chan error
	wg *sync.WaitGroup
}

func NewCrawler(output io.Writer) *Crawler {
	return &Crawler{
		output:            output,
		setOutputChannel:  make(chan []byte, 1000),
		visitedUris:       make(map[string]struct{}),
		setVisitedChannel: make(chan string, 1000),
		getVisitedChannel: make(chan uriRequest, 1000),
		urisToVisit:       make(chan string, 1000),
		stop:              make(chan error, 10),
		wg: new(sync.WaitGroup),
	}
}

func (c *Crawler) Run(uri string) (err error) {
	defer func() {
		e := recover();
		if err != nil {
			err = e.(error)
		}
	}()
	go func() {
		for uri := range c.urisToVisit {
			go c.Crawl(uri)
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func() {
		for {
			select {
			case uri := <-c.setVisitedChannel:
				c.visitedUris[uri] = struct{}{}
			case s := <-c.getVisitedChannel:
				_, ok := c.visitedUris[s.uri]
				s.respChannel <- ok
			case b := <-c.setOutputChannel:
				_, e := c.output.Write(b)
				if e != nil && e != io.EOF {
					c.stop <- e
				}
				c.output.Write([]byte("\n"))
			case err := <-c.stop:
				panic(err)
			}
		}
	}()
	c.ProcessUri(uri)
	c.wg.Wait()
	return err
}

func (c *Crawler) ProcessUri(uri string) {
	c.wg.Add(1)
	c.urisToVisit <- uri
}

func (c *Crawler) RememberAndOutput(uri string) {
	c.AddVisitedUri(uri)
	c.setOutputChannel <- []byte(uri)
}

func (c *Crawler) AddVisitedUri(uri string) {
	c.setVisitedChannel <- uri
}

func (c *Crawler) IsUriVisited(uri string) bool {
	rc := make(chan bool)
	c.getVisitedChannel <- uriRequest{
		uri:uri,
		respChannel:rc,
	}
	return <-rc
}

func (c *Crawler) Crawl(uri string) error {
	defer func() {
		c.wg.Done()
	}()
	if c.IsUriVisited(uri) {
		return nil
	}
	parsedUri, err := url.Parse(uri)
	if err != nil {
		return err
	}
	domain := []byte(parsedUri.Host + parsedUri.EscapedPath())
	r, _ := http.NewRequest(http.MethodGet, uri, nil)
	resp, err := (&http.Client{}).Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	c.RememberAndOutput(uri)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var p common.IParser
	p, err = parser.NewParser(nil, [][]byte{[]byte(domain)}, parsedUri)
	if err != nil {
		return err
	}
	//TODO fix fast parser
	//p = &parser.FastParser{}
	p.Init(b)
	for {
		u, ok := p.ParseNext();
		if !ok {
			break
		}
		c.ProcessUri(string(u))
	}
	return nil
}
