FROM golang:1.10-alpine as builder

RUN apk add --no-cache make curl git gcc musl-dev linux-headers
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ADD . /go/src/github.com/linkpoolio/clnode-exporter
RUN cd /go/src/github.com/linkpoolio/clnode-exporter && make build

# Copy stats exporter into a second stage container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/linkpoolio/clnode-exporter /usr/local/bin/

ADD docker/nodes.json /etc/exporter/nodes.json

EXPOSE 8080 8082
ENTRYPOINT ["clnode-exporter"]
CMD ["-configFile=/etc/exporter/nodes.json"]