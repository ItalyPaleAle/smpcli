## stkcli site set

Updates the configuration for a site

### Synopsis

Updates a site configured in the node.

When creating a site, you must specify the name of a TLS certificate stored in the node or cluster. Alternatively, you can pass one of the following values:

  - `selfsigned` for generating a self-signed certificate for your site
  - `acme` for requesting a certificate from an ACME provider, such as Let's Encrypt
  - `akv:[name]:[version]` for requesting a certificate stored in the Azure Key Vault instance associated with the cluster; the version is optional.

If you omit the `--certificate` option, it will default to a self-signed certificate.

The `--alias` parameter is used to replace the list of aliases configured for the domain. You can use this parameter multiple time to add more than one alias. Note that using the `--alias` flag will replace the entire list of aliases with the new one.


```
stkcli site set [flags]
```

### Options

```
  -a, --alias stringArray    alias domain (can be used multiple times)
  -c, --certificate string   name of the TLS certificate
  -d, --domain string        primary domain name
  -h, --help                 help for set
  -N, --node string          node address or IP
  -P, --port string          port the node listens on
```

### SEE ALSO

* [stkcli site](stkcli_site.md)	 - Manage sites

