FROM gcr.io/distroless/java

ADD spring /spring

ENTRYPOINT ["java", "org.springframework.boot.loader.PropertiesLauncher"]