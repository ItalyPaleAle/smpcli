## stkcli certificate add

Import a new TLS certificate

### Synopsis

Imports a new TLS certificate and stores it in the cluster's state.

You must provide a path to a PEM-encoded certificate and key using the `--certificate` and `--key` flags respectively.

The `--name` flag is the name of the TLS certificate used as identifier only.


```
stkcli certificate add [flags]
```

### Options

```
  -c, --certificate string   path to TLS certificate file
  -f, --force                force adding invalid/expired certificates
  -h, --help                 help for add
  -k, --key string           path to TLS key file
  -n, --name string          name for the certificate
  -N, --node string          node address or IP
  -P, --port string          port the node listens on
```

### SEE ALSO

* [stkcli certificate](stkcli_certificate.md)	 - Manage TLS certificates stored in the cluster

