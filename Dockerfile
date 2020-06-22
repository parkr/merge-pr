FROM golang:alpine AS builder
ARG REV
WORKDIR /usr/src/myapp
COPY . .
RUN set -ex \
  && apk add --no-cache git \
  && CI=1 CGO_ENABLED=0 time go test -v ./... \
  && CGO_ENABLED=0 go build -ldflags "-s -X main.revision=${REV}"

FROM golang:alpine
RUN apk add --no-cache git
COPY --from=builder /usr/src/myapp/merge-pr /bin/merge-pr
ENTRYPOINT [ "/bin/merge-pr" ]
