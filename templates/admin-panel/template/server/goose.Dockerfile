FROM golang:1.20.3-bullseye

RUN useradd -m goose && \
  go install github.com/pressly/goose/v3/cmd/goose@v3.18.0

CMD cd /server/migrations && \
    bash

