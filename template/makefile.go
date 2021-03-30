package template

var (
	escape   = "`grep \"asim\" -rl ./proto`"
	Makefile = `GOPATH:=$(shell go env GOPATH)
PROTOC_GEN=protoc --micro_out=paths=source_relative:. --gofast_out=paths=source_relative:. $$p

.PHONY: init
init:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get github.com/gogo/protobuf/protoc-gen-gofast
	go get github.com/asim/go-micro/cmd/protoc-gen-micro/v3

.PHONY: proto
proto:
	@for p in $(shell find proto -iname "*.proto"); do $(PROTOC_GEN); echo $(PROTOC_GEN); done
	@if [ $(shell uname) == "Darwin" ]; then sed -i "" s/"asim"/"blackdreamers"/g ` + escape + `; else sed -i s/"asim"/"blackdreamers"/g ` + escape + `; fi;
	
.PHONY: build
build:
	go build -ldflags "-s -w" -o {{.Alias}}-{{.Type}} *.go

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
