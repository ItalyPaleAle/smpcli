## stkcli upload

Upload apps and certificates

### Synopsis

The upload namespace contains commands to conveniently upload app bundles and TLS certificates.

IMPORTANT: In order to use these commands, you must have the Azure CLI installed and you must be authenticated to the Azure subscription where the Key Vault resides (with 'az login'). Additionally, your Azure account must have the following permissions in the Key Vault's data plane: keys (create, update, import, sign), certificate (create, update, import).


### Options

```
  -h, --help   help for upload
```

### SEE ALSO

* [stkcli](stkcli.md)	 - Manage a Statiko node
* [stkcli upload app](stkcli_upload_app.md)	 - Upload an app or bundle
* [stkcli upload certificate](stkcli_upload_certificate.md)	 - Upload a TLS certificate

