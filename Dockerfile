FROM golang:alpine AS builder

# TZ
RUN apk add --no-cache tzdata
ENV TZ=Asia/Tehran

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build --ldflags "-s -w -extldflags -static" -o main .

FROM alpine:latest

WORKDIR /www

COPY --from=builder /build/main /www/
COPY --from=builder /build/public/ /www/public/
COPY --from=builder /build/storage/ /www/storage/
COPY --from=builder /build/resources/ /www/resources/

# TZ
RUN apk add --no-cache tzdata
ENV TZ=Asia/Tehran

ENTRYPOINT ["/www/main"]
