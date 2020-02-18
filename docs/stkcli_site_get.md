## stkcli site get

Get a site

### Synopsis

Show the details of a site configured in the node.

Specify the primary domain name (no aliases) with the `--domain` parameter to select the site.


```
stkcli site get [flags]
```

### Options

```
  -d, --domain string   primary domain name
  -h, --help            help for get
  -S, --http            use HTTP protocol, without TLS, for node connections
  -k, --insecure        disable TLS certificate validation for node connections
  -n, --node string     node address or IP (required)
  -P, --port string     port the node listens on
```

### SEE ALSO

* [stkcli site](stkcli_site.md)	 - Manage sites

