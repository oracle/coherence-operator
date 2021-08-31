FROM scratch

ARG target

COPY lib/*.jar            /files/lib/
COPY logging/             /files/logging/
COPY linux/$target/runner /files/runner

ENTRYPOINT ["/files/runner"]
CMD ["-h"]
