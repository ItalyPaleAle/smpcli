## stkcli site add

Add a new site

### Synopsis

Configures a new site in the node.

Each site is identified by a primary domain, and it can have multiple aliases (domain names that are redirected to the primary one).

When creating a site, you can add the name of a TLS certificate stored on the associated Azure Key Vault instance. You can also specify 'selfsigned' as a value for the TLS certificate to have the node automatically generate a self-signed certificate for your site.


```
stkcli site add [flags]
```

### Options

```
  -a, --alias stringArray    alias domain (can be used multiple times)
  -c, --certificate string   name of the TLS certificate
  -d, --domain string        primary domain name
  -h, --help                 help for add
  -S, --http                 use HTTP protocol, without TLS
  -k, --insecure             disable TLS certificate validation (default true)
  -n, --node string          node address or IP (default "localhost")
  -P, --port string          port the node listens on (default "2265")
```

### SEE ALSO

* [stkcli site](stkcli_site.md)	 - Manage sites

