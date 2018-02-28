FROM golang:1-alpine

WORKDIR /go/src/

RUN apk add --no-cache git ca-certificates tzdata
RUN cp /usr/share/zoneinfo/Europe/Berlin /etc/localtime && echo "Europe/Berlin" >  /etc/timezone

RUN go get github.com/kabukky/klaus

CMD ["/go/bin/klaus"]
