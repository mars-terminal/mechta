FROM golang:1.21.0-alpine as builder

RUN apk add make

WORKDIR /builder

COPY . .

RUN make build-app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /builder/build/app ./app

ENTRYPOINT ["./app"]
