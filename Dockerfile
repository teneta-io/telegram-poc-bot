FROM golang:1.20-alpine3.17 AS builder

WORKDIR /app/

COPY . .

RUN go mod download
RUN set -x; apk add --no-cache \
    && CGO_ENABLED=0 go build -gcflags="all=-N -l"  \
    -a -installsuffix cgo -o ./bin/app cmd/bot/main.go

FROM alpine:3.17.3

WORKDIR /app

COPY --from=builder /app/bin .
COPY --from=builder /app/assets assets
COPY --from=builder /app/migrations migrations
COPY --from=builder /app/config.yml ./config.yml

RUN mkdir /app/secrets
RUN ln -s /vault/secrets/config.yml /app/config.yml

ENTRYPOINT ["./app"]