deps-build:
		go get -u github.com/golang/dep/...
		go get github.com/golang/lint/golint

deps: deps-build
		dep ensure

deps-update: deps-build
		rm -rf ./vendor
		rm -rf Gopkg.lock
		dep ensure -update

build:
		go build -ldflags="-w -s" -o bin/protodep main.go
