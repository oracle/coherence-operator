#
# Copyright (c) 2021 Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
ARG BASE_IMAGE
FROM $BASE_IMAGE

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT ["dlv"]
