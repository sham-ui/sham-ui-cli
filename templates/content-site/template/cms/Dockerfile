# Step 1: NodesJS modules cache
FROM node:14.21.3-bullseye as node_modules
RUN apt-get update && \
    apt-get install -y --no-install-recommends python2 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
COPY client/package.json client/yarn.lock /client/
WORKDIR /client
RUN yarn install --frozen-lockfile --ignore-scripts --verbose

# Step 2: Build the js files
FROM node:14.21.3-bullseye as build_js
COPY --from=node_modules /client/node_modules /client/node_modules
COPY client/ /client
WORKDIR /client
RUN yarn build

# Step 3: Go modules & dependencies cache
FROM golang:1.20.3-bullseye as go_modules
RUN apt-get update && \
    apt-get install -y --no-install-recommends make bash git gcc libc-dev tzdata nodejs npm && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*\ && \
    npm install -g n && \
    n 18.14.0 && \
    wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.51.2
COPY server/go.mod server/go.sum /server/
WORKDIR /server
RUN go mod download

# Step 4: Run lint, test, build & run integration tests
FROM go_modules as build_go
# Copy directoies
COPY --from=build_js client/dist/ /server/assets/files/
# Copy files
COPY --from=build_js client/dist/* /server/assets/files/
COPY server/ /server
WORKDIR /server
RUN echo "Run linter" && golangci-lint run --timeout=15m ./... && \
    echo "Run go test" && go test -p 1 -v ./cmd/... ./config/... ./internal/... ./pkg/...  && \
    echo "Make build" && GOOS=linux GOARCH=amd64 go build -x -o /bin/app ./cmd/app && \
    echo "Run integration tests" && cp /bin/app ./integration-test/app && go test -p 1 -v ./integration-test/... && \
    echo "Success"

# Step 5: Final build
FROM debian:bullseye-slim as final
WORKDIR /
COPY --from=build_go /bin/app .
ENV MALLOC_ARENA_MAX=4
ENTRYPOINT ["./app"]



