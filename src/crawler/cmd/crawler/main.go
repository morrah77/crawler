package main

import (
	"flag"
	"io"
	"os"
	"log"
	"strings"
	"crawler/crawler"
)

type Conf struct {
	url string
	output string
}

var conf Conf

func init() {
	flag.StringVar(&conf.url, `url`, ``, `Url to crawl, for ex., http://google.com`)
	flag.StringVar(&conf.output, `output`, ``, `filename to write output (std out if empty)`)
	flag.Parse()
}

func main() {
	var (
		writer io.Writer
		err error
		cr *crawler.Crawler
	)
	if strings.Trim(conf.output, ` `) == `` {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(conf.output, os.O_CREATE, os.ModePerm)
		failOnError(err)
	}
	cr = crawler.NewCrawler(writer)
	println(conf.url)
	failOnError(cr.Run(conf.url))
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
