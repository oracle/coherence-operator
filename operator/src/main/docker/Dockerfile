FROM openjdk:11-oracle

RUN mkdir /operator

COPY * /operator/
RUN chmod +x /operator/*.sh
COPY lib/* /operator/lib/

WORKDIR /operator/

ENTRYPOINT ["./run.sh"]
