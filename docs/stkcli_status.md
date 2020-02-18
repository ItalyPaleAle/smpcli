## stkcli status

Shows the status of a node

### Synopsis

Prints information about the status and health of the node.

The '--domain' flag allows selecting a specific site only.


```
stkcli status [flags]
```

### Options

```
  -d, --domain string   domain name
  -h, --help            help for status
  -S, --http            use HTTP protocol, without TLS
  -k, --insecure        disable TLS certificate validation (default true)
  -n, --node string     node address or IP (default "localhost")
  -P, --port string     port the node listens on (default "2265")
```

### SEE ALSO

* [stkcli](stkcli.md)	 - Manage a Statiko node

