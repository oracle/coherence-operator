# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2025, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
# tag::dockerfile[]

FROM node:18-alpine

WORKDIR /usr/src/app

COPY package*.json ./
COPY main.js ./

RUN npm install

EXPOSE 8080

CMD [ "node", "main.js" ]

# end::dockerfile[]