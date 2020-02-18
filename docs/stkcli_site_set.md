## stkcli site set

Updates the configuration for a site

### Synopsis

Updates a site configured in the node.

Use the '--certificate' parameter to set a new TLS certificate. This should be the name of a certificate stored in the associated Azure Key Vault. You can also use the value "selfsigned" to have the node automatically generate a self-signed certificate for your site.

The '--alias' parameter is used to replace the list of aliases configured for the domain. You can use this parameter multiple time to add more than one alias. Note that using the '--alias' flag will replace the entire list of aliases with the new one.


```
stkcli site set [flags]
```

### Options

```
  -a, --alias stringArray    alias domain (can be used multiple times)
  -c, --certificate string   name of the TLS certificate
  -d, --domain string        primary domain name
  -h, --help                 help for set
  -S, --http                 use HTTP protocol, without TLS
  -k, --insecure             disable TLS certificate validation (default true)
  -n, --node string          node address or IP (default "localhost")
  -P, --port string          port the node listens on (default "2265")
```

### SEE ALSO

* [stkcli site](stkcli_site.md)	 - Manage sites

