FROM golang:1.10.3

WORKDIR /go/src/github.com/HcashOrg/hcd
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go install . ./cmd/...

EXPOSE 14008

CMD hcd
