#The simplest web crawler

##Original description

Задание:
Нужно написать паука, который принимает на вход url веб страницы. Паук определяет домен этого url и ищет страницы этого домена по ссылкам, доступных на страницах этого домена. На выходе имеем список урлов домена. 


Требования:


* Паук должен беречь сайты, для этого ограничим кол-во запросов в секунду;
* Найденные урлы страниц просто пишутся в консоль;
* Исходный код лежит на github;
* Паук должен использовать все доступные ресурсы; 
* Выбор языка программирования за кандидатом, желательно обоснование;
* Паук запускается через docker. Например сл. образом: docker build test-crawler . && docker run test-crawler https://some-url.com/some-page.

##Description

###REM ONLY PLAIN HTML IS PARSED! IS NOT PURPOSED TO PARSE GENERATED CONTENT!

 :Crawl
 - increment Crawl threads counter
 - if url still not visited
   - get response from url
   - if got (status code < 500)
     - mark url as visited
     - if status is OK
       - write url
       - parse response
       - for each item found
         - Crawl item in concurrent thread
 - Decrement Crawl threads counter

 - wait until all Crawl threads return


So, crawler should have thread-safe visited URIs storage, channel of URI to be visited, thread counter and run each child Crawl function in separate goroutine not faster than specified times a second.


cr := New Crawler()
cr.Run()

Crawler:

Run(uri) {
  go func() {
    for uri := range c.urisToVisit {
      go Crawl(uri)
      time.Sleep(10 * time.Millisecond)
    }
  }()
  c.urisToVisit <- uri
  for {
    if len(c.urisToVisit) == 0 && c.counter == 0 {
      return
    }
  }
}

##Build

###Locally

```
export GOPATH=`pwd` &&\
cd src/crawler &&\
dep ensure &&\
go install ./cmd/... &&\
cd - &&\
```

###Docker

`docker build -t crawler:latest -f ./Dockerfile .`

##Test

```
cd src/crawler/
go test ./...

```

###Locally

```
export GOPATH=`pwd` &\
go test crawler/...
```

###Docker

##Run

###Locally

`./bin/crawler -url https://morrah77.ru`

###Docker

`docker run -it -e "URL=https://morrah77.ru" crawler:latest`

or

```
mkdir -p output &\
docker run --rm -v `pwd`/output:/usr/local/crawler crawler -e "URL=https://morrah77.ru" -e "OUTPUT="/usr/local/crawler/file.txt" &\
tail -f -n 1000 `pwd`/output/file.txt
```
