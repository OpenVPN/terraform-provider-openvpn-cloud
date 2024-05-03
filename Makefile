HOSTNAME=cloudconnexa.dev
NAMESPACE=openvpn
NAME=cloudconnexa
VERSION=0.0.12
BINARY=terraform-provider-${NAME}
OS_ARCH=darwin_arm64

default: install

build:
	go build -o ${BINARY}

release:
	goreleaser release --clean --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

lint:
	golangci-lint run ./...

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
