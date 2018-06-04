build:
	go build ./cmd/manssh

build-docker:
	docker build --build-arg VERSION=`git describe --tags` -t wendellsun/manssh .
