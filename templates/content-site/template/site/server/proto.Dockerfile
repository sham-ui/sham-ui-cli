FROM golang:1.20.3-bullseye

RUN cd / && \
    apt-get update && \
    apt-get install -y --no-install-recommends unzip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v24.4/protoc-24.4-linux-x86_64.zip && \
    unzip -o protoc-24.4-linux-x86_64.zip -d /usr/local/ bin/protoc && \
    unzip -o protoc-24.4-linux-x86_64.zip -d /usr/local/ 'include/*' && \
    rm protoc-24.4-linux-x86_64.zip && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

CMD cd /site && \
    find . -name '*.proto' | xargs -I {} protoc -I=. -I=${GOPATH}/pkg/mod/ --experimental_allow_proto3_optional --go-grpc_out=../ --go_out=../ {}

