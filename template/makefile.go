package template

var (
	Makefile = `GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto
PROTOC_GEN=protoc --micro_out=paths=source_relative:. --gofast_out=paths=source_relative:. $$p

.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/micro/go-micro/cmd/protoc-gen-micro/v2

.PHONY: proto
proto:
	@for p in $(shell find proto -iname "*.proto"); do $(PROTOC_GEN); echo $(PROTOC_GEN); done
	
.PHONY: build
build:
	go build -o {{.Alias}}-{{.Type}} *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t {{.Alias}}:latest
`

	GenerateFile = `package main
//go:generate make proto
`
)
