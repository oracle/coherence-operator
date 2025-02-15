FROM ghcr.io/oracle/coherence-operator:3.4.3 AS Builder

FROM ghcr.io/oracle/coherence-ce:22.06.10
COPY --from=Builder /files /files
COPY --from=Builder /files/lib/coherence-operator.jar /app/libs/coherence-operator.jar
COPY coherence-java-client-22.06.7.jar /app/libs/coherence-java-client-22.06.7.jar

ENTRYPOINT ["files/runner"]
CMD ["-h"]