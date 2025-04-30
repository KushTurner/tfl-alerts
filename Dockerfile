FROM golang:1.24.2 AS builder

ARG CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build-deploy

FROM alpine:latest

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]