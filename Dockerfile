FROM golang:1.5.3

WORKDIR /go/src/github.com/anarcher/glia
ADD . /go/src/github.com/anarcher/glia
ENV GOPATH /go/src/github.com/anarcher/glia/Godeps/_workspace:$GOPATH
RUN go install -v

ENTRYPOINT ["glia"]
