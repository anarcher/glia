FROM golang:1.5.3

WORKDIR /go/src/github.com/anarcher/glia
ADD . /go/src/github.com/anarcher/glia
ENV GOPATH /go/src/github.com/anarcher/glia/Godeps/_workspace:$GOPATH
RUN ./.build_version
RUN go install -ldflags="-X github.com/anarcher/glia/Version=`cat ./VERSION` -X github.com/anarcher/glia/GitCommit=`git rev-parse --short HEAD`" -v 

ENTRYPOINT ["glia"]
