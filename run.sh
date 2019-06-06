#!/bin/bash

URL={$1:-https://morrah77.ru}

OUTPUT=$2

docker build -f ./Dockerfile .

docker run --rm -it crawler ${URL:+"-e \"URL=$URL\""} ${OUTPUT:+"-e \"OUTPUT=OUTPUT\""}
