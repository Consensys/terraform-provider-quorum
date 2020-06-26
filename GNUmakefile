TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=quorum
OUTPUT_DIR=$(shell readlink build)
CURRENT_OS=$(shell go env GOOS)
CURRENT_ARCH=$(shell go env GOARCH)
VERSION=$(shell cat VERSION)
TARGET_OS=linux
ifeq (darwin,$(CURRENT_OS))
	TARGET_OS+=local
endif
PLATFORMS=$(addprefix dist, $(TARGET_OS))
GOLANG_DOCKER_IMAGE=golang:1.13-alpine
LDFLAGS=-s -w $(ldflags)

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

dist: fmtcheck $(PLATFORMS)

distlinux:
	@echo "==> Building for Linux using Docker with $(GOLANG_DOCKER_IMAGE)"
	@docker run --rm -t -v $(shell pwd):/terraform-provider-quorum -v $(OUTPUT_DIR):$(OUTPUT_DIR) $(GOLANG_DOCKER_IMAGE) /terraform-provider-quorum/scripts/linux_build.sh

distlocal:
	@echo '==> Building for $(CURRENT_OS)_$(CURRENT_ARCH) with -ldflags: $(LDFLAGS)'
	@GOFLAGS="-mod=vendor" go build -o build/$(CURRENT_OS)_$(CURRENT_ARCH)/terraform-provider-quorum_v$(VERSION) -ldflags '$(LDFLAGS)'

vendor:
	@go mod vendor
	@cp -rf $(shell go list -f {{.Dir}} github.com/ethereum/go-ethereum/crypto/secp256k1) ./vendor/github.com/ethereum/go-ethereum/crypto/
	@chmod -R u+w ./vendor/github.com/ethereum/go-ethereum/crypto/
	@cp -rf $(shell go list -f {{.Dir}} github.com/karalabe/hid)/ ./vendor/github.com/karalabe/hid/
	@chmod -R u+w ./vendor/github.com/karalabe/hid/

docs:
	@cd website && go generate

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test dist distlinux distdarwin vendor