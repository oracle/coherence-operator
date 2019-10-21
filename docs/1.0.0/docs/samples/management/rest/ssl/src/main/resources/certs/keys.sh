#!/usr/bin/env bash

if [ "${COMPUTERNAME}" == "" ]
then
    COMPUTERNAME="localhost"
fi

if [ "${CLIENTCOMPUTERNAME}" == "" ]
then
    CLIENTCOMPUTERNAME="localhost"
fi

echo Generate Guardians CA key:
openssl genrsa -passout pass:1111 -des3 -out guardians-ca.key 4096

echo Generate Guardians CA certificate:
openssl req -passin pass:1111 -new -x509 -days 3650 -key guardians-ca.key -out guardians-ca.crt -subj "/CN=${COMPUTERNAME}" # guardians-ca.crt is a trustCertCollectionFile

echo Generate Ravagers CA key:
openssl genrsa -passout pass:1111 -des3 -out ravagers-ca.key 4096

echo Generate Ravagers CA certificate:
openssl req -passin pass:1111 -new -x509 -days 3650 -key ravagers-ca.key -out ravagers-ca.crt -subj "/CN=${COMPUTERNAME}" # ravagers-ca.crt is a trustCertCollectionFile


echo Generate Icarus key:
openssl genrsa -passout pass:1111 -des3 -out icarus.key 4096

echo Generate Icarus signing request:
openssl req -passin pass:1111 -new -key icarus.key -out icarus.csr -subj "/CN=${COMPUTERNAME}"

echo Self-signed Icarus certificate:
openssl x509 -req -passin pass:1111 -days 3650 -in icarus.csr -CA guardians-ca.crt -CAkey guardians-ca.key -set_serial 01 -out icarus.crt # icarus.crt is the certChainFile for the server

echo Remove passphrase from Icarus key:
openssl rsa -passin pass:1111 -in icarus.key -out icarus.key


echo Generate client Star-Lord key
openssl genrsa -passout pass:1111 -des3 -out star-lord.key 4096

echo Generate client Star-Lord signing request:
openssl req -passin pass:1111 -new -key star-lord.key -out star-lord.csr -subj "/CN=Star-Lord"

echo Self-signed client Star-Lord certificate:
openssl x509 -passin pass:1111 -req -days 3650 -in star-lord.csr -CA guardians-ca.crt -CAkey guardians-ca.key -set_serial 01 -out star-lord.crt # star-lord.crt is the certChainFile for the client (Mutual TLS only)

echo Remove passphrase from Star-Lord key:
openssl rsa -passin pass:1111 -in star-lord.key -out star-lord.key


echo Generate client Groot key
openssl genrsa -passout pass:1111 -des3 -out groot.key 4096

echo Generate client Groot signing request:
openssl req -passin pass:1111 -new -key groot.key -out groot.csr -subj "/CN=Groot"

echo Self-signed client Groot certificate:
openssl x509 -passin pass:1111 -req -days 3650 -in groot.csr -CA guardians-ca.crt -CAkey guardians-ca.key -set_serial 01 -out groot.crt # groot.crt is the certChainFile for the client (Mutual TLS only)

echo Remove passphrase from client Groot key:
openssl rsa -passin pass:1111 -in groot.key -out groot.key


echo Generate client Yondu key
openssl genrsa -passout pass:1111 -des3 -out yondu.key 4096

echo Generate client Yondu signing request:
openssl req -passin pass:1111 -new -key yondu.key -out yondu.csr -subj "/CN=Yondu"

echo Self-signed client Yondu certificate:
openssl x509 -passin pass:1111 -req -days 3650 -in yondu.csr -CA ravagers-ca.crt -CAkey ravagers-ca.key -set_serial 01 -out yondu.crt # yondu.crt is the certChainFile for the client (Mutual TLS only)

echo Remove passphrase from client Yondu key:
openssl rsa -passin pass:1111 -in yondu.key -out yondu.key


openssl pkcs8 -topk8 -nocrypt -in star-lord.key -out star-lord.pem # star-lord.pem is the privateKey for the Client (mutual TLS only)
openssl pkcs8 -topk8 -nocrypt -in groot.key -out groot.pem # groot.pem is the privateKey for the Client (mutual TLS only)
openssl pkcs8 -topk8 -nocrypt -in yondu.key -out yondu.pem # yondu.pem is the privateKey for the Client (mutual TLS only)
openssl pkcs8 -topk8 -nocrypt -in icarus.key -out icarus.pem # icarus.pem is the privateKey for the Server

# Create the Java trust store
rm *.jks

keytool -import -storepass secret -noprompt -trustcacerts -alias guardians -file guardians-ca.crt -keystore truststore-guardians.jks -deststoretype JKS
keytool -import -storepass secret -noprompt -trustcacerts -alias ravagers -file ravagers-ca.crt -keystore truststore-ravagers.jks -deststoretype JKS
keytool -import -storepass secret -noprompt -trustcacerts -alias guardians -file guardians-ca.crt -keystore truststore-all.jks -deststoretype JKS
keytool -import -storepass secret -noprompt -trustcacerts -alias ravagers -file ravagers-ca.crt -keystore truststore-all.jks -deststoretype JKS

openssl pkcs12 -export -passout pass:password -inkey icarus.pem -name test -in icarus.crt -out icarus.p12
keytool -importkeystore -storepass password -noprompt -srcstorepass password -srckeystore icarus.p12 -srcstoretype pkcs12 -destkeypass password -destkeystore icarus.jks

openssl pkcs12 -export -passout pass:password -inkey star-lord.pem -name test -in star-lord.crt -out star-lord.p12
keytool -importkeystore -storepass password -noprompt -srcstorepass password -srckeystore star-lord.p12 -srcstoretype pkcs12 -destkeypass password -destkeystore star-lord.jks

openssl pkcs12 -export -passout pass:password -inkey groot.pem -name test -in groot.crt -out groot.p12
keytool -importkeystore -storepass password -noprompt -srcstorepass password -srckeystore groot.p12 -srcstoretype pkcs12 -destkeypass password -destkeystore groot.jks

openssl pkcs12 -export -passout pass:password -inkey yondu.pem -name test -in yondu.crt -out yondu.p12
keytool -importkeystore -storepass password -noprompt -srcstorepass password -srckeystore yondu.p12 -srcstoretype pkcs12 -destkeypass password -destkeystore yondu.jks

rm *.crt
rm *.csr
rm *.key
rm *.p12
rm *.pem
