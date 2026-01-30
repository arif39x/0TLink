#!/bin/bash

mkdir -p certs

#the Certificate Authority (CA)
# This is the 'Root of Trust'. Both server and client will trust this.
openssl genrsa -out certs/ca.key 2048
openssl req -new -x509 -days 365 -key certs/ca.key -out certs/ca.crt -subj "/CN=0TLink-CA"

#the Server Certificate (The Relay)
openssl genrsa -out certs/server.key 2048
openssl req -new -key certs/server.key -out certs/server.csr -subj "/CN=localhost"
# Sign the server cert with our CA
openssl x509 -req -days 365 -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt

#the Client Certificate (The Agent)
openssl genrsa -out certs/client.key 2048
openssl req -new -key certs/client.key -out certs/client.csr -subj "/CN=0TLink-Agent"
# Sign the client cert with our CA
openssl x509 -req -days 365 -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/client.crt

# Clean up temporary signing requests
rm certs/*.csr certs/*.srl

echo "Success: Zero-Trust certs generated in ./certs"