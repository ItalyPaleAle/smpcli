## stkcli dhparams set

Sets new DH parameters

### Synopsis

Sets new Diffie-Hellman parameters for the cluster. If the cluster is currently re-generating them, this interrupts the operation.

The --`file` flag is the path to a PEM-encoded file containing DH parameters.


```
stkcli dhparams set [flags]
```

### Options

```
  -f, --file string   path to DH parameters file
  -h, --help          help for set
  -N, --node string   node address or IP
  -P, --port string   port the node listens on
```

### SEE ALSO

* [stkcli dhparams](stkcli_dhparams.md)	 - Set DH parameters for the cluster

