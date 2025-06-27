FROM golang:1.23-alpine as builder
RUN apk --no-cache --no-progress add make git
WORKDIR /go/tlsctl
ENV GO111MODULE on
COPY go.mod .
COPY go.sum .
#RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
RUN make build
FROM alpine:3
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=builder /go/tlsctl/tlsctl /usr/bin/tlsctl
ENTRYPOINT [ "/usr/bin/tlsctl" ]