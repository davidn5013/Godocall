## Example makefile for Go with fuzz test

## This is for not using go.mod
# .EXPORT_ALL_VARIABLES:
# GOPATH=$(shell cd)
# GO111MODULE=off

default: run

help:
	@echo Usage:
	@echo make help  - this help
	@echo make qc    - go vet 
	@echo make lint  - All linter and error checkors 
	@echo make run   - run with go vet
	@echo make test  - run tests 
	@echo make bench - run bench
	@echo make build - build with full lint and staticcheck
	@echo make clean - delete *.exe

qc: 
	@-go vet .
	@echo -----------------------------------------------------------------

lint: qc
	-golangci-lint run .  
	-revive . 
	-errcheck . 
	@-echo -----------------------------------------------------------------

test: qc
	-go test ./...
	@echo -----------------------------------------------------------------

bench: qc 
	go test -benchmem -run=. -bench=. -benchtime=20s
	@echo -----------------------------------------------------------------
	go test -fuzz=./... -fuzztime=20s 
	@echo -----------------------------------------------------------------
	
build: lint test 
	go build -gcflags="-m=2" --ldflags="-s -w -race" -trimpath .
	@echo -----------------------------------------------------------------

run: qc
	go run .
	@echo -----------------------------------------------------------------

.PHONY : clean
clean :
	@go clean .
	@-rm *.exe      
