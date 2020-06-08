## stkcli app upload

Upload an app or bundle

### Synopsis

Uploads an app or app bundle to the node, to be stored in the node's app repository.

This command accepts four parameters:

- `--path` is the path to a file or folder to upload
- `--app` is the name of the name of the bundle, which can be used to identify the app when you want to deploy it in a node (do not include an extension)
- `--signing-key` is the path to a private RSA key used for codesigning

Paths can be folders containing your app's files; stkcli will automatically create a tar.bz2 archive for you. Alternatively, you can point the `--path` parameter to an existing archive (various formats are supported, including zip, tar.gz, tar.bz2, and more), and it will uploaded as-is.

App names must be unique. You cannot re-upload an app using the same file name.


```
stkcli app upload [flags]
```

### Options

```
  -a, --app string           app bundle name, with no extension (required)
  -h, --help                 help for upload
  -N, --node string          node address or IP
  -f, --path string          path to local file or folder to bundle (required)
  -P, --port string          port the node listens on
  -s, --signing-key string   path to a RSA private key for code signing
```

### SEE ALSO

* [stkcli app](stkcli_app.md)	 - Upload and manage app bundles

