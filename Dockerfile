FROM golang:1.21.6-alpine as builder

ARG LD_FLAGS="-s -w"
ARG TARGETPLATFORM
ARG TARGETARCH

WORKDIR /app

COPY server/pkg pkg
COPY server/main.go main.go
COPY server/go.sum go.sum
COPY server/go.mod go.mod

RUN export GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) && \
    export GOARCH=$(TARGETARCH) && \
    apk --no-cache add ca-certificates && \
    update-ca-certificates

RUN go env

RUN go get -d -v \
    && go install -v

RUN CGO_ENABLED=0 go build -ldflags="${LD_FLAGS}" -o /app/build/policyreporter-ui -v

FROM scratch
LABEL MAINTAINER "Frank Jogeleit <frank.jogeleit@gweb.de>"

WORKDIR /app

USER 1234

COPY LICENSE.md .
COPY dist dist
COPY --from=builder /app/build/policyreporter-ui /app/policyreporter-ui
# copy the debian's trusted root CA's to the final image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENTRYPOINT ["/app/policyreporter-ui"]
