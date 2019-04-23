FROM golang:latest as builder
WORKDIR /go/src/github.com/jsimonetti/go-artnet
ADD . /go/src/github.com/jsimonetti/go-artnet
ENV CGO_ENABLED=0
RUN go build -ldflags '-w -extldflags "-static"' -o node example/node/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /go/src/github.com/jsimonetti/go-artnet/node /app/
CMD ["/app/node"]