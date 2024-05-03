GOPATH:=$(shell go env GOPATH)


.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@go build -o server cmd/server/main.go

.PHONY: docker
docker:
	@docker build -t server:latest .