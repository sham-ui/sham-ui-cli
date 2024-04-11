GOOSE_DEV?=postgres://{{ dbUser }}:{{ dbPassword }}@localhost:5432/{{ dbName }}
GOOSE_TEST?=postgres://{{ dbUser }}:{{ dbPassword }}@localhost:5432/{{ dbName }}_test

install:
	go mod download -x

goose_dev:
	docker run --rm -h={{ shortName }}-goose-dev -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING=${GOOSE_DEV} --net=host -u $(shell id -u ${USER}):$(shell id -g ${USER}) -v $(shell pwd):/server -it $(shell docker build -t {{ shortName }}-goose -f goose.Dockerfile -q .)

goose_test_up:
	docker run --rm -h={{ shortName }}-goose-test -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING=${GOOSE_TEST} --net=host -u $(shell id -u ${USER}):$(shell id -g ${USER}) -v $(shell pwd):/server -it $(shell docker build -t {{ shortName }}-goose -f goose.Dockerfile -q .) goose -dir /server/migrations up

.ONESHELL:
build_js:
	cd ../client/
	yarn build
	rm -r ../server/assets/files
	mkdir -p ../server/assets/files
	cp -r dist/* ../server/assets/files

build_amd64: build_js
	GOOS=linux GOARCH=amd64 go build -x -o ./bin/linux_amd64/app ./cmd/app

build: build_amd64

run_test:
	go test -p 1 -v ./cmd/... ./config/... ./internal/... ./migrations/... ./pkg/... -cover -coverprofile=coverage.out -coverpkg=./... && echo "Tests finished with success"

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
