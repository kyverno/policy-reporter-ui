FROM node:15-alpine as frontend

WORKDIR /var/www

RUN apk add --no-cache --virtual .gyp make g++ \
    && npm set progress=false \
    && npm config set depth 0

COPY frontend frontend

RUN cd frontend \
    && npm install \
    && npm run build

FROM golang:1.17.2 as builder

ARG LD_FLAGS
ARG TARGETPLATFORM

WORKDIR /app

COPY pkg pkg
COPY main.go main.go
COPY go.sum go.sum
COPY go.mod go.mod
COPY Makefile Makefile

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2)

RUN go env

RUN go get -d -v \
    && go install -v

RUN CGO_ENABLED=0 go build -ldflags="${LD_FLAGS}" -o /app/build/policyreporter-ui -v

FROM scratch
LABEL MAINTAINER "Frank Jogeleit <frank.jogeleit@gweb.de>"

WORKDIR /app

USER 1234

COPY LICENSE.md .
COPY --from=frontend /var/www/dist /app/dist
COPY --from=builder /app/build/policyreporter-ui /app/policyreporter-ui

EXPOSE 2112

ENTRYPOINT ["/app/policyreporter-ui"]
