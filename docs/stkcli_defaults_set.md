## stkcli defaults set

Set default connection options

### Synopsis

Set the value for the shared flags that will be used as default in all commands:

- '--node address' or '-n address':
  Sets the address (IP or hostname) of the node to connect to.
  This option is required.
- '--port port' or '-P port':
  If set, will communicate with the node using the port specified.
  System default: 2265
- '--insecure' or '-k' (boolean):
  If set, disables TLS certificate validation when communicating with the node (e.g. to use self-signed certificates).
  System default: false (requires valid TLS certificate)
- '--http' or '-S' (boolean):
  If set, communicates with the node using unencrypted HTTP.
  This option is considered insecure, and should only be used if the node is 'localhost', or if you're connecting to the node over an already-encrypted tunnel (e.g. VPN or SSH port forwarding).
  System default: false (use TLS)

Note that calling the 'defaults set' command overrides the default values for all the four flags above. If those values are not set, the system defaults are used. 


```
stkcli defaults set [flags]
```

### Options

```
  -h, --help          help for set
  -S, --http          use HTTP protocol, without TLS
  -k, --insecure      disable TLS certificate validation (default true)
  -n, --node string   node address or IP (default "localhost")
  -P, --port string   port the node listens on (default "2265")
```

### SEE ALSO

* [stkcli defaults](stkcli_defaults.md)	 - View or set defaults for stkcli

