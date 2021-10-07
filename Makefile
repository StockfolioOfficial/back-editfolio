ifndef BINARY
	BINARY=debug
endif

init:
	go mod download all
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	go install github.com/swaggo/swag/cmd/swag
	go install github.com/google/wire/cmd/wire

generate-env:
	export PATH=$PATH:$GOPATH

generate: generate-env swagger wire

swagger:
	swag init

wire:
	wire .

proto-compile:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

test:
	go test -v -cover -covermode=atomic ./...

build:
	go build -o ${BINARY} .

unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

lint-prepare:
	@echo "Installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...