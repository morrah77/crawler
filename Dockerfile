FROM golang:1.12

ENV URL=https://morrah77.ru
ENV OUTPUT=
WORKDIR /go
COPY src src
RUN go get -u github.com/golang/dep/cmd/dep
RUN cd src/crawler && dep ensure && cd -
#RUN go test crawler/...
RUN go install crawler/cmd/...
CMD ./bin/crawler -url $URL ${OUTPUT:+-output $OUTPUT}
