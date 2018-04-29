FROM golang:1.10 as builder
ARG VERSION
WORKDIR /go/src/github.com/xwjdsh/manssh
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION}" -o manssh ./cmd/manssh

FROM alpine:latest  
LABEL maintainer="iwendellsun@gmail.com"
WORKDIR /
COPY --from=builder /go/src/github.com/xwjdsh/manssh/manssh .
ENTRYPOINT ["/manssh"]
