## stkcli auth psk

Authenticate using a pre-shared key

### Synopsis

Sets the pre-shared key used to authenticate API calls to a node.

The pre-shared key is defined in the node's configuration, and clients are authenticated if they send the same key in the header of API calls.
Note that the key is not hashed nor encrypted, so using TLS to connect to nodes is strongly recommended.


```
stkcli auth psk [flags]
```

### Options

```
  -h, --help          help for psk
  -S, --http          use HTTP protocol, without TLS
  -k, --insecure      disable TLS certificate validation
  -n, --node string   node address or IP (required)
  -P, --port string   port the node listens on
```

### SEE ALSO

* [stkcli auth](stkcli_auth.md)	 - Authenticate with a node

