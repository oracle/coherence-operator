# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2025, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
# tag::dockerfile[]

FROM python:3.11-slim

RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

WORKDIR /app
COPY --chown=appuser:appgroup main.py .
RUN chmod 444 main.py

RUN pip install --no-cache-dir coherence-client==2.0.0 Quart

RUN chown -R appuser:appgroup /app

USER appuser

CMD ["python3", "./main.py"]

# end::dockerfile[]