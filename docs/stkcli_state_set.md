## stkcli state set

Restores the state of a node

### Synopsis

Replaces the state of the node with the one read from file (or stdin if the '--file' parameter is not set).

The state is a JSON document containing the list of sites and apps to be configured in the web server, and it's normally exported from another node (useful for backups or migrations).

This command completely replaces the state of the node with the one you're passing to the command, discarding any site or app currently configured in the node.


```
stkcli state set [flags]
```

### Options

```
  -f, --file string   file containing the desired state; if not set, read from stdin
  -h, --help          help for set
  -S, --http          use HTTP protocol, without TLS
  -k, --insecure      disable TLS certificate validation
  -n, --node string   node address or IP (required)
  -P, --port string   port the node listens on
```

### SEE ALSO

* [stkcli state](stkcli_state.md)	 - Get or restore state

