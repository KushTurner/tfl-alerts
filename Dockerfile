FROM golang:1.23.1 AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .

RUN go build -o main .

FROM scratch
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]