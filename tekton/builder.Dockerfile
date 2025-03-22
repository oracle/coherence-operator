FROM container-registry.oracle.com/os/oraclelinux:9

ARG GoVersion=1.24.1
ARG GoArch=arm64

CMD ["/bin/bash"]

RUN dnf install make which -y

RUN curl -L https://go.dev/dl/go$GoVersion.linux-$GoArch.tar.gz -o go-linux.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go-linux.tar.gz \
    && rm go-linux.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"
