# docker build -t "customer-logs-sample1.0.0-SNAPSHOT" .

FROM oraclelinux:7-slim
RUN mkdir -p /files/lib /files/conf
COPY files/lib/*  files/lib/
COPY files/conf/* files/conf/
