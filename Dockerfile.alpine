# Build image
FROM golang:1.10.3

WORKDIR /go/src/github.com/HcashOrg/hcd
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go install . ./cmd/...

# Production image
FROM alpine:3.6

RUN apk add --no-cache ca-certificates
COPY --from=0 /go/bin/* /bin/

EXPOSE 14008

CMD hcd