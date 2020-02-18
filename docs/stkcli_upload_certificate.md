## stkcli upload certificate

Upload a TLS certificate

### Synopsis

Adds a TLS certificate (public certificate and private key) to the Azure Key Vault instance connected to the node.

IMPORTANT: In order to use this command, you must have the Azure CLI installed and you must be authenticated to the Azure subscription where the Key Vault resides (with 'az login'). Additionally, your Azure account must have the following permissions in the Key Vault's data plane: keys (create, update, import, sign), certificate (create, update, import).

This command accepts three parameters:

- '--name' or '-c' is the name of the certificate, which you can use to reference it in sites' configuration. Per Azure, it can only contain uppercase and lowercase letters, numbers and the dash symbol (-)
- '--certificate' or '-f' is the file with the public part of the TLS certificate
- '--certificate-key' or '-p' is the file with the private key of the certificate

Note that only certificates with RSA keys are supported. Additionally, both the certificate and the key must be in PEM format.


```
stkcli upload certificate [flags]
```

### Options

```
  -f, --certificate string       certificate file (required)
  -p, --certificate-key string   private key (required)
  -h, --help                     help for certificate
  -S, --http                     use HTTP protocol, without TLS
  -k, --insecure                 disable TLS certificate validation (default true)
  -c, --name string              certificate name (required)
  -n, --node string              node address or IP (default "localhost")
  -P, --port string              port the node listens on (default "2265")
```

### SEE ALSO

* [stkcli upload](stkcli_upload.md)	 - Upload apps and certificates

