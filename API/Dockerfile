FROM golang:1.23.8 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main .

FROM build-stage AS test-stage
RUN go test -v ./...

FROM debian:bookworm-slim AS production-stage

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build-stage /app/main .

RUN chmod +x main

USER nobody:nogroup

ENTRYPOINT ["/app/main"]