FROM golang:1.21-alpine

RUN apk update && apk upgrade && apk add --update alpine-sdk && \
  apk add --no-cache bash make cmake 

WORKDIR /funcgo

COPY . .

RUN make docker

CMD [ "./bin/controller" ]
