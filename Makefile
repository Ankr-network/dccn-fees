.PHONY: build test run deploy local clean

GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools

build:
	docker build . -t app_dccn_dcmgr:latest

run: build
	docker run -d \
		$(MICRO_ENV) \
		$(PROGRAM_ENV) \
		app_dccn_dcmgr

local:
	$(MICRO_ENV) \
	$(PROGRAM_ENV) \
	go run main.go

clean:
	rm app_dccn_dcmgr

test:
	go test -v ./... -cover -race

deploy:
	@echo "docker push"

devtools:
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go get github.com/micro/protoc-gen-micro
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
	@type "protoc-gen-micro" 2> /dev/null || echo 'Please install protoc-gen-micro'


define MICRO_ENV
    MICRO_SERVER_ADDRESS=":50051" \
	MICRO_BROKER_ADDRESS="amqp://guest:guest@localhost:5672"\
	DEV_EVN=false
endef

define PROGRAM_ENV
	DB_HOST="127.0.0.1:27017" \
	DB_NAME="dccn"
endef

