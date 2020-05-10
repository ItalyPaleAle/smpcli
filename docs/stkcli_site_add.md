## stkcli site add

Add a new site

### Synopsis

Configures a new site in the node.

Each site is identified by a primary domain, and it can have multiple aliases (domain names that are redirected to the primary one).

When creating a site, you must specify the name of a TLS certificate stored on the associated Azure Key Vault instance. You can also specify `selfsigned` as a value for the TLS certificate to have the node automatically generate a self-signed certificate for your site. If you omit the `--certificate` option, it will default to a self-signed certificate.


```
stkcli site add [flags]
```

### Options

```
  -a, --alias stringArray        alias domain (can be used multiple times)
  -c, --certificate selfsigned   name of the TLS certificate or selfsigned (default)
  -d, --domain string            primary domain name
  -h, --help                     help for add
  -n, --node string              node address or IP
  -P, --port string              port the node listens on
```

### SEE ALSO

* [stkcli site](stkcli_site.md)	 - Manage sites

