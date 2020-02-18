## stkcli site remove

Remove a site

### Synopsis

Removes a site from the node, so the web server stops accepting requests for it.

You must specify the primary domain name (no aliases) in the `--domain` parameter to select the site to be removed.


```
stkcli site remove [flags]
```

### Options

```
  -d, --domain string   primary domain name
  -h, --help            help for remove
  -S, --http            use HTTP protocol, without TLS
  -k, --insecure        disable TLS certificate validation
  -n, --node string     node address or IP (required)
  -P, --port string     port the node listens on
      --yes             do not ask for confirmation
```

### SEE ALSO

* [stkcli site](stkcli_site.md)	 - Manage sites

