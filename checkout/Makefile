CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PROTODIR=${CURDIR}/../api
PROTOC=protoc
PACKAGE=route256/checkout/cmd/app

all: format build test lint

build: bindir build-proto
	go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks -exclude pkg

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

install-grpc-bins = \
	${BINDIR}/protoc-gen-go \
	${BINDIR}/protoc-gen-go-grpc \
	${BINDIR}/protoc-gen-grpc-gateway \
	${BINDIR}/protoc-gen-openapiv2

$(install-grpc-bins):
	GOBIN=${BINDIR} go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=${BINDIR} go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=${BINDIR} go get -u github.com/grpc-ecosystem/grpc-gateway/v2

	GOBIN=${BINDIR} go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
	GOBIN=${BINDIR} go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=${BINDIR} go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=${BINDIR} go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2

install-grpc: $(install-grpc-bins)

${PROTODIR}/googleapis:
	git clone https://github.com/googleapis/googleapis.git ${PROTODIR}/googleapis

build-proto: install-grpc ${PROTODIR}/googleapis $(shell find ${PROTODIR} -name "*.proto" -type f)
	mkdir -p pkg/checkout
	${PROTOC} -I ${PROTODIR}/checkout -I ${PROTODIR}/googleapis \
		--go_out=pkg/checkout --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/checkout --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		--grpc-gateway_out=pkg/checkout --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
		${PROTODIR}/checkout/checkout.proto

	mkdir -p pkg/loms
	${PROTOC} -I ${PROTODIR}/loms -I ${PROTODIR}/googleapis \
		--go_out=pkg/loms --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/loms --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		${PROTODIR}/loms/loms.proto

	mkdir -p pkg/product
	${PROTOC} -I ${PROTODIR}/product -I ${PROTODIR}/googleapis \
		--go_out=pkg/product --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/product --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		${PROTODIR}/product/product.proto