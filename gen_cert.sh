#!/bin/bash
mkdir -p certs
# CA
openssl genrsa -out certs/ca.key 2048
openssl req -new -x509 -days 365 -key certs/ca.key -out certs/ca.crt -subj "/CN=0TLink-CA"
# Server
openssl genrsa -out certs/server.key 2048
openssl req -new -key certs/server.key -out certs/server.csr -subj "/CN=localhost"
openssl x509 -req -days 365 -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial \
    -out certs/server.crt -extfile <(printf "subjectAltName=DNS:localhost,IP:127.0.0.1")
# Client
openssl genrsa -out certs/client.key 2048
openssl req -new -key certs/client.key -out certs/client.csr -subj "/CN=0TLink-Agent"
openssl x509 -req -days 365 -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial \
    -out certs/client.crt -extfile <(printf "subjectAltName=DNS:0TLink-Agent")
rm certs/*.csr certs/*.srl
