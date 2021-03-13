FROM node:14-alpine as frontend

WORKDIR /var/www

RUN apk add --no-cache --virtual .gyp make g++ \
    && npm set progress=false \
    && npm config set depth 0

COPY frontend frontend

RUN cd frontend \
    && npm install \
    && npm run build

FROM golang:1.16-buster as builder

WORKDIR /app
COPY pkg pkg
COPY main.go main.go
COPY go.sum go.sum
COPY go.mod go.mod
COPY Makefile Makefile

RUN go get -d -v \
    && go install -v

RUN make build

FROM alpine:latest
LABEL MAINTAINER "Frank Jogeleit <frank.jogeleit@gweb.de>"

WORKDIR /app

RUN apk add --update --no-cache ca-certificates

RUN addgroup -S policyreporter && adduser -u 1234 -S policyreporter -G policyreporter

USER 1234

COPY LICENSE.md .
COPY --from=frontend /var/www/dist /app/dist
COPY --from=builder /app/build/policyreporter-ui /app/policyreporter-ui

EXPOSE 2112

ENTRYPOINT ["/app/policyreporter-ui"]