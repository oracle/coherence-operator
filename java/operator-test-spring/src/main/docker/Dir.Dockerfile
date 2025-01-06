FROM gcr.io/distroless/java17-debian12

ADD spring /spring

ENTRYPOINT ["java", "org.springframework.boot.loader.PropertiesLauncher"]