.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

APP := protodep
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

build: tidy vendor
		GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-w -s" -o bin/protodep -mod=vendor main.go

define build-artifact
		GO111MODULE=on GOOS=$(1) GOARCH=$(2) go build -ldflags="-w -s" -o artifacts/$(3) -mod=vendor main.go
		cd artifacts && tar cvzf $(APP)_$(1)_$(2).tar.gz $(3)
		rm ./artifacts/$(3)
		@echo [INFO]build success: $(1)_$(2)
endef

.PHONY: build-all
build-all: tidy vendor
		$(call build-artifact,linux,386,$(APP))
		$(call build-artifact,linux,amd64,$(APP))
		$(call build-artifact,linux,arm,$(APP))
		$(call build-artifact,linux,arm64,$(APP))
		$(call build-artifact,darwin,amd64,$(APP))
		$(call build-artifact,windows,amd64,$(APP).exe)

.PHONY: clean
clean:
	rm -rf bin
	rm -rf vendor
	rm -rf artifacts 

.PHONY: test
test:
	go test -v ./...
