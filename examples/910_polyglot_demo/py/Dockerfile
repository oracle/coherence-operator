# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2025, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
# tag::dockerfile[]

FROM python:3.11-slim

ADD main.py .
RUN pip install coherence-client==2.0.0 Quart

CMD ["python3", "./main.py"]

# end::dockerfile[]