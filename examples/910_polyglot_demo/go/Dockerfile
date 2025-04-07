# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2025, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
# tag::dockerfile[]

FROM golang:1.24 as builder

# Set working dir inside the container
WORKDIR /app

# Copy your Go module files and source code
COPY go.mod ./
COPY go.sum ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o runner .

FROM scratch

COPY --from=builder /app/runner /files/runner

EXPOSE 8080
ENTRYPOINT ["/files/runner"]
CMD ["-h"]

# end::dockerfile[]