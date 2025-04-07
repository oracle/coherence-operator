# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2025, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
# tag::dockerfile[]

FROM node:18-alpine

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /usr/src/app

COPY --chown=appuser:appgroup package*.json ./
COPY --chown=appuser:appgroup main.js ./

RUN chmod 444 package*.json main.js

RUN npm install --ignore-scripts
RUN chown -R appuser:appgroup /usr/src/app

USER appuser
EXPOSE 8080

CMD [ "node", "main.js" ]

# end::dockerfile[]