ROOT_PACKAGE := github.com/stormcat24/protodep
VERSION_PACKAGE := $(ROOT_PACKAGE)/version
LDFLAG_GIT_COMMIT := "$(VERSION_PACKAGE).gitCommit"
LDFLAG_GIT_COMMIT_FULL := "$(VERSION_PACKAGE).gitCommitFull"
LDFLAG_BUILD_DATE := "$(VERSION_PACKAGE).buildDate"
LDFLAG_VERSION := "$(VERSION_PACKAGE).version"

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

APP := protodep
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

build: tidy vendor version
		$(eval GIT_COMMIT := $(shell git describe --tags --always))
		$(eval GIT_COMMIT_FULL := $(shell git rev-parse HEAD))
		$(eval BUILD_DATE := $(shell date '+%Y%m%d'))
		GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-w -s -X $(LDFLAG_GIT_COMMIT)=$(GIT_COMMIT) -X $(LDFLAG_GIT_COMMIT_FULL)=$(GIT_COMMIT_FULL) -X $(LDFLAG_BUILD_DATE)=$(BUILD_DATE) -X $(LDFLAG_VERSION)=$(GIT_COMMIT)" \
			-o bin/protodep -mod=vendor main.go

define build-artifact
		$(eval GIT_COMMIT := $(shell git describe --tags --always))
		$(eval GIT_COMMIT_FULL := $(shell git rev-parse HEAD))
		$(eval BUILD_DATE := $(shell date '+%Y%m%d'))
		GO111MODULE=on GOOS=$(1) GOARCH=$(2) go build -ldflags="-w -s -X $(LDFLAG_GIT_COMMIT)=$(GIT_COMMIT) -X $(LDFLAG_GIT_COMMIT_FULL)=$(GIT_COMMIT_FULL) -X $(LDFLAG_BUILD_DATE)=$(BUILD_DATE) -X $(LDFLAG_VERSION)=$(BUILD_DATE)-$(GIT_COMMIT)" \
			-o artifacts/$(3) -mod=vendor main.go
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
		$(call build-artifact,darwin,arm64,$(APP))
		$(call build-artifact,windows,amd64,$(APP).exe)

.PHONY: clean
clean:
	rm -rf bin
	rm -rf vendor
	rm -rf artifacts 

.PHONY: test
test:
	go test -v ./...

.PHONY: update-credits
update-credits:
	@go install github.com/Songmu/gocredits/cmd/gocredits@latest
	@gocredits . > CREDITS
