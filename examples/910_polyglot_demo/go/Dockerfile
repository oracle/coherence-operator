# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2025, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
# tag::dockerfile[]

FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY main.go ./

ENV APP_USER_UID=1000
ENV APP_USER_GID=1000

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o runner .
RUN chown ${APP_USER_UID}:${APP_USER_GID} /app/runner

FROM scratch

COPY --from=builder /app/runner /files/runner
USER 1000:1000

EXPOSE 8080
ENTRYPOINT ["/files/runner"]
CMD ["-h"]

# end::dockerfile[]