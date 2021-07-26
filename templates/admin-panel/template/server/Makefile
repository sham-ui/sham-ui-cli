restore_deps:
	cd {{ shortName }}
	go mod download

install: restore_deps

.ONESHELL:
build_js:
	cd ../client/
	yarn build
	mkdir -p ../server/{{ shortName }}/assets/files
	cp -r dist/* ../server/{{ shortName }}/assets/files

build_386: build_js
	cd {{ shortName }}
	env GOOS=linux GOARCH=386 go build -o ../bin/linux_386/{{ shortName }} {{ shortName }}

build_amd64: build_js
	cd {{ shortName }}
	env GOOS=linux GOARCH=amd64 go build -o ../bin/linux_amd64/{{ shortName }} {{ shortName }}

build: build_386 build_amd64

run_test:
	cd {{ shortName }}
	go test -v ./... -cover -coverprofile=coverage.out -coverpkg=./... && echo "Tests finished with success"

coverage:
	cd {{ shortName }}
	go tool cover -html=coverage.out

test: build_js run_test coverage
