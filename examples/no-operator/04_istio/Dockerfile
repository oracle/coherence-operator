FROM container-registry.oracle.com/middlewarecoherence-operator:3.5.4 AS Builder

FROM container-registry.oracle.com/middleware/coherence-ce:14.1.2-0-1
COPY --from=Builder /files /files
COPY --from=Builder /files/lib/coherence-operator.jar /app/libs/coherence-operator.jar
COPY coherence-java-client-14.1.2-0-1.jar /app/libs/coherence-java-client-14.1.2-0-1.jar

ENTRYPOINT ["files/runner"]
CMD ["-h"]