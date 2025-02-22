FROM gcr.io/distroless/java17-debian12

COPY spring /spring

WORKDIR /spring

ENTRYPOINT ["java", "org.springframework.boot.loader.launch.PropertiesLauncher"]