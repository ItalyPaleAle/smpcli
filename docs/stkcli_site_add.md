## stkcli site add

Add a new site

### Synopsis

Configures a new site in the node.

Each site is identified by a primary domain, and it can have multiple aliases (domain names that are redirected to the primary one).

Alternatively, you can specify the `--temporary` option to create a temporary site, for example for testing an application. When creating temporary sites, a domain name will be generated for you, and you should not provide domain names or aliases.

When creating a site, you must specify the name of a TLS certificate stored in the node or cluster. Alternatively, you can pass one of the following values:

  - `selfsigned` for generating a self-signed certificate for your site
  - `acme` for requesting a certificate from an ACME provider, such as Let's Encrypt (not available for temporary sites)
  - `akv:[name]:[version]` for requesting a certificate stored in the Azure Key Vault instance associated with the cluster; the version is optional.

If you omit the `--certificate` option, it will default to a self-signed certificate.


```
stkcli site add [flags]
```

### Options

```
  -a, --alias stringArray        alias domain (can be used multiple times)
  -c, --certificate selfsigned   name of the TLS certificate or selfsigned (default)
  -d, --domain string            primary domain name (required for non-temporary sites)
  -h, --help                     help for add
  -n, --node string              node address or IP
  -P, --port string              port the node listens on
  -t, --temporary                create a temporary site with a random name
```

### SEE ALSO

* [stkcli site](stkcli_site.md)	 - Manage sites

