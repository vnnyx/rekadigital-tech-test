FROM golang:1.20-alpine as builder
WORKDIR /builder
COPY . .
RUN apk add --no-cache upx \
    && go mod download \
    && go build -ldflags "-s -w" -o main \
    && upx -9 main

FROM alpine:latest
WORKDIR /app
COPY --from=builder /builder/main .
COPY --from=builder /builder/configs ./configs
COPY --from=builder /builder/migrations/ ./migrations
CMD ["/app/main", "server"]