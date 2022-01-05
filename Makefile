build:
	go build ./cmd/manssh

docker-build:
	docker build --build-arg VERSION=`git describe --tags` -t wendellsun/manssh .

docker-push:
	docker push wendellsun/manssh
