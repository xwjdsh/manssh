FROM golang:1.17 as builder
ARG VERSION
WORKDIR /go/src/github.com/xwjdsh/manssh
COPY . .
RUN go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION}" ./cmd/manssh

FROM alpine:3.15
LABEL maintainer="iwendellsun@gmail.com"
WORKDIR /
COPY --from=builder /go/src/github.com/xwjdsh/manssh/manssh .
ENTRYPOINT ["/manssh"]
