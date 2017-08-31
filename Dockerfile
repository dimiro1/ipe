FROM golang:latest 
#ADD . /go/src/github.com/dimiro1/ipe
RUN go get github.com/dimiro1/ipe

ENV PORT 8080
EXPOSE 8080
ENTRYPOINT /go/bin/ipe -config /app/config.json -logtostderr -v 0
