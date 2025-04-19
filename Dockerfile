FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache make
RUN make build

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/bin/core ./bin/core
COPY configs/config.yaml ./configs/config.yaml
COPY migrations migrations

CMD ["./bin/core", "./configs/config.yaml"]