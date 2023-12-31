| Environmental Variable | Default | Description |
| -- | -- | -- |
| TLS_CERT_FILE | none | Cert file location for TLS Cert  |
| TLS_KEY_FILE | none | Key file location for TLS Cert |
| PORT | `8080` | Port for HTTP Server running inside container |
| FILES_DIR | `/public` | Port for browsing directory inside container |

# Running SpectroCloud Browser
```
docker run -p 8080:8080 -v /shared:/app/public spectrocloud-browser
```

# SpectroCloud Browser w/ TLS

## Generating Local Test Certificates

```
openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
```

```
docker run -p 443:443 \
    -v /shared:/app/public \
    -v localhost.crt:/localhost.crt \
    -v localhost.key:/localhost.key \
    -e TLS_CERT_FILE=/localhost.crt \
    -e TLS_KEY_FILE=/localhost.key \
    spectrocloud-browser
```
