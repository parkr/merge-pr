FROM golang:alpine AS builder
WORKDIR /usr/src/myapp
COPY . .
RUN go build -ldflags '-s'

FROM golang:alpine
RUN apk add --no-cache git
COPY --from=builder /usr/src/myapp/merge-pr /bin/merge-pr
ENTRYPOINT [ "/bin/merge-pr" ]
