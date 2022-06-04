FROM node:16-alpine as frontend

WORKDIR /var/www

RUN apk add --no-cache --virtual .gyp python3 make g++ \
    && npm set progress=false \
    && npm config set depth 0

COPY . .

RUN npm install \
    && npm run generate

FROM golang:1.18 as builder

ARG LD_FLAGS
ARG TARGETPLATFORM

WORKDIR /app

COPY server/pkg pkg
COPY server/main.go main.go
COPY server/go.sum go.sum
COPY server/go.mod go.mod

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

EXPOSE 8080

ENTRYPOINT ["/app/policyreporter-ui"]
