## stkcli upload app

Upload an app or bundle

### Synopsis

Uploads an app to the Azure Storage Account associated with the node's instance

IMPORTANT: In order to use this command, you must have the Azure CLI installed and you must be authenticated to the Azure subscription where the Key Vault resides (with `az login`). Additionally, your Azure account must have the following permissions in the Key Vault's data plane: keys (create, update, import, sign), certificate (create, update, import).

This command accepts four parameters:

- `--path` or `-f` is the path of the app to upload
- `--app` or `-a` is the name of the app's bundle, which can be used to identify the app when you want to deploy it in a node
- `--version` or `-v` is the version of the app, as an arbitrary string
- `--no-signature` is a boolean that when present will skip calculating the checksum of the app's bundle and signing it with the codesign key

Paths can be folders containing your app's files; stkcli will automatically create a tar.bz2 archive for you. Alternatively, you can point the `--path` parameter to an existing tar.bz2 archive, and it will uploaded as-is.

Versions are unique for each app. For example, if you upload the app `myapp` and version `1.0`, you cannot re-upload that; the version must be different. Statiko does not parse the version and does not enforce any specific versioning convention, as long as the versions are different.

When using `--no-signature`, stkcli will not calculate the checksum of the app's bundle, and it will not cryptographically sign it with the codesigning key. Statiko nodes might be configured to not accept unsigned app bundles for security reasons. However, when uploading unsigned bundles, you do not need to be signed into an Azure account in the local system.


```
stkcli upload app [flags]
```

### Options

```
  -a, --app string       app's bundle name (required)
  -h, --help             help for app
      --no-signature     do not cryptographically sign the app's bundle
  -n, --node string      node address or IP (required)
  -f, --path string      path to local file or folder to bundle
  -P, --port string      port the node listens on
  -v, --version string   app's bundle version (required)
```

### SEE ALSO

* [stkcli upload](stkcli_upload.md)	 - Upload apps and certificates

