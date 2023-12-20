restore_deps:
	go mod download -x

install: restore_deps

build_proto:
	docker run --rm -u $(shell id -u ${USER}):$(shell id -g ${USER}) -v $(shell pwd):/site -it $(shell docker build -t site-proto -f proto.Dockerfile -q .)

.ONESHELL:
build_js:
	cd ../client/
	yarn build
	rm -rf ../server/assets/files/
	mkdir -p ../server/assets/files
	cp -r dist/* ../server/assets/files
	cp ssr/bundle.js ../server/ssr/ssr.js

build_amd64: build_proto build_js
	GOOS=linux GOARCH=amd64 go build -x -o ./bin/linux_amd64/app ./cmd/app

build: build_amd64

run_test:
	go test -p 1 -v ./cmd/... ./config/... ./internal/... ./pkg/...  -cover -coverprofile=coverage.out -coverpkg=./... && echo "Tests finished with success"

coverage:
	go tool cover -html=coverage.out

test: build_js run_test coverage

run_integration_test:
	go test -p 1 -v ./integration-test/...  && echo "Integrations tests finished with success"

cp_app_to_integration_test:
	cp bin/linux_amd64/app integration-test/app

integration_test: build cp_app_to_integration_test run_integration_test

lint:
	golangci-lint run ./...

lint_fix:
	golangci-lint run --fix ./...

run_amd64:
	./bin/linux_amd64/app

run: run_amd64
