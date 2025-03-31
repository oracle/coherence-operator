FROM container-registry.oracle.com/os/oraclelinux:9

ARG GoVersion
ARG GoArch

CMD ["/bin/bash"]

RUN dnf install oracle-java-jdk-release-el* -y \
    && dnf install jdk-21-headful -y \
    && dnf install make which git -y

RUN curl -Ls https://go.dev/dl/$GoVersion.linux-$GoArch.tar.gz -o go-linux.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go-linux.tar.gz \
    && rm go-linux.tar.gz \
    && curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/$GoArch/kubectl" \
    && chmod u+x kubectl \
    && mv kubectl /usr/local/bin

ENV PATH="$PATH:/usr/local/go/bin"
