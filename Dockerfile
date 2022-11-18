FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/gitlab.com/tunder-tunder/avito/
WORKDIR /go/src/gitlab.com/tunder-tunder/avito/
RUN go mod download
COPY . /go/src/gitlab.com/tunder-tunder/avito/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/avito gitlab.com/idoko/avito

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/tunder-tunder/avito/build/avito /usr/bin/avito
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/avito"]