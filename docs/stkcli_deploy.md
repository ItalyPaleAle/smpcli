## stkcli deploy

Deploy an app

### Synopsis

Deploys an app to a site.

This command tells the node to deploy the app (already uploaded beforehand) with the specific name and version to a site identified by the domain option.


```
stkcli deploy [flags]
```

### Options

```
  -a, --app string       app's bundle name (required)
  -d, --domain string    primary domain name (required)
  -h, --help             help for deploy
  -n, --node string      node address or IP (required)
  -P, --port string      port the node listens on
  -v, --version string   app's bundle version (required)
```

### SEE ALSO

* [stkcli](stkcli.md)	 - Manage a Statiko node

