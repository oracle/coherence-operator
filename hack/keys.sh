#!/usr/bin/env bash

if [ "${COMPUTERNAME}" == "" ]
then
    COMPUTERNAME="localhost"
fi

if [ "${CLIENTCOMPUTERNAME}" == "" ]
then
    CLIENTCOMPUTERNAME="localhost"
fi

mkdir -p build/_output/certs

echo Generate Guardians CA key:
openssl genrsa -passout pass:1111 -aes256 \
    -out build/_output/certs/guardians-ca.key 4096

echo Generate Guardians CA certificate:
openssl req -passin pass:1111 -new -x509 -days 3650 \
    -key build/_output/certs/guardians-ca.key \
    -out build/_output/certs/guardians-ca.crt \
    -subj "/CN=${COMPUTERNAME}" # guardians-ca.crt is a trustCertCollectionFile

echo Generate Ravagers CA key:
openssl genrsa -passout pass:1111 -aes256 \
    -out build/_output/certs/ravagers-ca.key 4096

echo Generate Ravagers CA certificate:
openssl req -passin pass:1111 -new -x509 -days 3650 \
    -key build/_output/certs/ravagers-ca.key \
    -out build/_output/certs/ravagers-ca.crt \
    -subj "/CN=${COMPUTERNAME}" # ravagers-ca.crt is a trustCertCollectionFile


echo Generate Icarus key:
openssl genrsa -passout pass:1111 -aes256 \
    -out build/_output/certs/icarus.key 4096

echo Generate Icarus signing request:
openssl req -passin pass:1111 -new -key \
    build/_output/certs/icarus.key \
    -out build/_output/certs/icarus.csr \
    -subj "/CN=${COMPUTERNAME}"

echo Self-signed Icarus certificate:
openssl x509 -req -passin pass:1111 -days 3650 \
    -in build/_output/certs/icarus.csr \
    -CA build/_output/certs/guardians-ca.crt \
    -CAkey build/_output/certs/guardians-ca.key \
    -set_serial 01 \
    -out build/_output/certs/icarus.crt # icarus.crt is the certChainFile for the server

echo Remove passphrase from Icarus key:
openssl rsa -passin pass:1111 \
    -in build/_output/certs/icarus.key \
    -out build/_output/certs/icarus.key


echo Generate client Star-Lord key
openssl genrsa -passout pass:1111 -aes256 \
    -out build/_output/certs/star-lord.key 4096

echo Generate client Star-Lord signing request:
openssl req -passin pass:1111 -new \
    -key build/_output/certs/star-lord.key \
    -out build/_output/certs/star-lord.csr -subj "/CN=Star-Lord"

echo Self-signed client Star-Lord certificate:
openssl x509 -passin pass:1111 -req -days 3650 \
    -in build/_output/certs/star-lord.csr \
    -CA build/_output/certs/guardians-ca.crt \
    -CAkey build/_output/certs/guardians-ca.key \
    -set_serial 01 \
    -out build/_output/certs/star-lord.crt # star-lord.crt is the certChainFile for the client (Mutual TLS only)

echo Remove passphrase from Star-Lord key:
openssl rsa -passin pass:1111 \
    -in build/_output/certs/star-lord.key \
    -out build/_output/certs/star-lord.key

echo Generate client Groot key
openssl genrsa -passout pass:1111 -aes256 \
    -out build/_output/certs/groot.key 4096

echo Generate client Groot signing request:
openssl req -passin pass:1111 -new \
    -key build/_output/certs/groot.key \
    -out build/_output/certs/groot.csr \
    -subj "/CN=Groot"

echo Self-signed client Groot certificate:
openssl x509 -passin pass:1111 -req -days 3650 \
    -in build/_output/certs/groot.csr -CA build/_output/certs/guardians-ca.crt \
    -CAkey build/_output/certs/guardians-ca.key \
    -set_serial 01 \
    -out build/_output/certs/groot.crt # groot.crt is the certChainFile for the client (Mutual TLS only)

echo Remove passphrase from client Groot key:
openssl rsa -passin pass:1111 \
    -in build/_output/certs/groot.key \
    -out build/_output/certs/groot.key


echo Generate client Yondu key
openssl genrsa -passout pass:1111 -aes256 \
    -out build/_output/certs/yondu.key 4096

echo Generate client Yondu signing request:
openssl req -passin pass:1111 -new \
    -key build/_output/certs/yondu.key \
    -out build/_output/certs/yondu.csr \
    -subj "/CN=Yondu"

echo Self-signed client Yondu certificate:
openssl x509 -passin pass:1111 -req -days 3650 \
    -in build/_output/certs/yondu.csr \
    -CA build/_output/certs/ravagers-ca.crt \
    -CAkey build/_output/certs/ravagers-ca.key \
    -set_serial 01 \
    -out build/_output/certs/yondu.crt # yondu.crt is the certChainFile for the client (Mutual TLS only)

echo Remove passphrase from client Yondu key:
openssl rsa -passin pass:1111 \
    -in build/_output/certs/yondu.key \
    -out build/_output/certs/yondu.key


openssl pkcs8 -topk8 -nocrypt \
    -in build/_output/certs/star-lord.key \
    -out build/_output/certs/star-lord.pem # star-lord.pem is the privateKey for the Client (mutual TLS only)

openssl pkcs8 -topk8 -nocrypt \
    -in build/_output/certs/groot.key \
    -out build/_output/certs/groot.pem # groot.pem is the privateKey for the Client (mutual TLS only)

openssl pkcs8 -topk8 -nocrypt \
    -in build/_output/certs/yondu.key \
    -out build/_output/certs/yondu.pem # yondu.pem is the privateKey for the Client (mutual TLS only)

openssl pkcs8 -topk8 -nocrypt \
    -in build/_output/certs/icarus.key \
    -out build/_output/certs/icarus.pem # icarus.pem is the privateKey for the Server

# Create the Java trust store
rm build/_output/certs/*.jks

KEYPASS="p455w0rd"
STOREPASS="p455w0rd"
TRUSTPASS="secret"

keytool -import -storepass ${TRUSTPASS} -noprompt -trustcacerts \
    -alias guardians -file build/_output/certs/guardians-ca.crt \
    -keystore build/_output/certs/truststore-guardians.jks \
    -deststoretype JKS

keytool -import -storepass ${TRUSTPASS} -noprompt -trustcacerts \
    -alias ravagers -file build/_output/certs/ravagers-ca.crt \
    -keystore build/_output/certs/truststore-ravagers.jks \
    -deststoretype JKS

keytool -import -storepass ${TRUSTPASS} -noprompt -trustcacerts \
    -alias guardians -file build/_output/certs/guardians-ca.crt \
    -keystore build/_output/certs/truststore-all.jks \
    -deststoretype JKS

keytool -import -storepass ${TRUSTPASS} -noprompt -trustcacerts \
    -alias ravagers -file build/_output/certs/ravagers-ca.crt \
    -keystore build/_output/certs/truststore-all.jks \
    -deststoretype JKS

openssl pkcs12 -export -passout pass:${KEYPASS} \
    -inkey build/_output/certs/icarus.pem \
    -name test -in build/_output/certs/icarus.crt \
    -out build/_output/certs/icarus.p12

keytool -importkeystore -storepass ${STOREPASS} -noprompt \
    -srcstorepass ${KEYPASS} \
    -srckeystore build/_output/certs/icarus.p12 \
    -srcstoretype pkcs12 \
    -destkeypass ${KEYPASS} \
    -destkeystore build/_output/certs/icarus.jks

openssl pkcs12 -export -passout pass:${KEYPASS} \
    -inkey build/_output/certs/star-lord.pem \
    -name test -in build/_output/certs/star-lord.crt \
    -out build/_output/certs/star-lord.p12

keytool -importkeystore -storepass ${STOREPASS} -noprompt \
    -srcstorepass ${KEYPASS} \
    -srckeystore build/_output/certs/star-lord.p12 \
    -srcstoretype pkcs12 \
    -destkeypass ${KEYPASS} \
    -destkeystore build/_output/certs/star-lord.jks

openssl pkcs12 -export -passout pass:${KEYPASS} \
    -inkey build/_output/certs/groot.pem \
    -name test -in build/_output/certs/groot.crt \
    -out build/_output/certs/groot.p12

keytool -importkeystore -storepass ${STOREPASS} -noprompt \
    -srcstorepass ${KEYPASS} \
    -srckeystore build/_output/certs/groot.p12 \
    -srcstoretype pkcs12 \
    -destkeypass ${KEYPASS} \
    -destkeystore build/_output/certs/groot.jks

openssl pkcs12 -export -passout pass:${KEYPASS} \
    -inkey build/_output/certs/yondu.pem \
    -name test -in build/_output/certs/yondu.crt \
    -out build/_output/certs/yondu.p12

keytool -importkeystore -storepass ${STOREPASS} -noprompt \
    -srcstorepass ${KEYPASS} \
    -srckeystore build/_output/certs/yondu.p12 \
    -srcstoretype pkcs12 \
    -destkeypass ${KEYPASS} \
    -destkeystore build/_output/certs/yondu.jks

echo ${KEYPASS} > build/_output/certs/keypassword.txt
echo ${KEYPASS} > build/_output/certs/storepassword.txt
echo ${TRUSTPASS} > build/_output/certs/trustpassword.txt
