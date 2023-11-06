TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=CiscoDevNet
NAME=secureworkload
BINARY=terraform-provider-${NAME}
VERSION=0.2
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: install

build:
	go build -o ${BINARY}_${VERSION}_${OS_ARCH}

PLATFORMS := linux/amd64 windows/amd64 darwin/amd64 darwin/arm64

temp = $(subst /, ,$@)
OS = $(word 1, $(temp))
ARCH = $(word 2, $(temp))

release: $(PLATFORMS)
	tar cvzf release.tgz release
	zip -r release.zip release

$(PLATFORMS):
	mkdir -p ./release/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${ARCH}
	GOOS=${OS} GOARCH=${ARCH} go build -o ./release/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${ARCH}/${BINARY}_${VERSION}_${OS}_${ARCH}

.PHONY: release $(PLATFORMS)

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY}_${VERSION}_${OS_ARCH} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

generate:
	tfplugindocs
