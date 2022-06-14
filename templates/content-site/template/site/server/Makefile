restore_deps:
	cd site
	go mod download

install: restore_deps

build_proto:
	cd site
	find ./proto -name '*.proto' | xargs -I {} protoc -I=. -I=${GOPATH}/pkg/mod/ --go-grpc_out=../ --go_out=../ {}

.ONESHELL:
build_js:
	cd ../client/
	yarn build
	rm -r ../server/site/assets/files
	mkdir -p ../server/site/assets/files
	cp -r dist/* ../server/site/assets/files
	cp ssr/bundle.js ../server/site/ssr/ssr.js

build_386: build_js
	cd site
	env GOOS=linux GOARCH=386 go build -o ../bin/linux_386/site site

build_amd64: build_js
	cd site
	env GOOS=linux GOARCH=amd64 go build -o ../bin/linux_amd64/site site

build: build_386 build_amd64

run_test:
	cd site
	go test -p 1 -v ./... -cover -coverprofile=coverage.out -coverpkg=./... && echo "Tests finished with success"

coverage:
	cd site
	go tool cover -html=coverage.out

test: build_js run_test coverage
