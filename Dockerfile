FROM golang:1.23.1 AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .

RUN make build-deploy

FROM alpine:latest
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]