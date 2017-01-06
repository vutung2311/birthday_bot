FROM golang:alpine

WORKDIR /go/src/birthday-bot

RUN apk add git build-base sqlite-dev && \
    go get github.com/golang/dep/cmd/dep

COPY . /go/src/birthday-bot

RUN dep ensure -v

CMD go run /go/src/birthday-bot/main.go