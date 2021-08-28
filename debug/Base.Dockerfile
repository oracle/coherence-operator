#
# Copyright (c) 2021 Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
FROM golang:1.16

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

COPY debug/debug.sh .

RUN go get github.com/go-delve/delve/cmd/dlv
RUN chmod +x debug.sh

ENTRYPOINT ["dlv"]
