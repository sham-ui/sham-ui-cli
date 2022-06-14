restore_deps:
	cd cms
	go mod download

install: restore_deps

build_proto:
	cd cms
	find ./proto -name '*.proto' | xargs -I {} protoc -I=. -I=${GOPATH}/pkg/mod/ --go-grpc_out=../ --go_out=../ {}

.ONESHELL:
build_js:
	cd ../client/
	yarn build
	rm -r ../server/cms/assets/files
	mkdir -p ../server/cms/assets/files
	cp -r dist/* ../server/cms/assets/files

build_386: build_js
	cd cms
	env GOOS=linux GOARCH=386 go build -o ../bin/linux_386/cms cms

build_amd64: build_js
	cd cms
	env GOOS=linux GOARCH=amd64 go build -o ../bin/linux_amd64/cms cms

build: build_386 build_amd64

run_test:
	cd cms
	go test -p 1 -v ./... -cover -coverprofile=coverage.out -coverpkg=./... && echo "Tests finished with success"

coverage:
	cd cms
	go tool cover -html=coverage.out

test: build_js run_test coverage
