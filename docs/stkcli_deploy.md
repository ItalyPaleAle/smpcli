## stkcli deploy

Deploy an app

### Synopsis

Deploys an app to a site.

This command tells the node to deploy the app (already uploaded beforehand) with the specific bundle name to a site identified by the domain option.


```
stkcli deploy [flags]
```

### Options

```
  -a, --app string      app bundle (required)
  -d, --domain string   primary domain name (required)
  -h, --help            help for deploy
  -N, --node string     node address or IP
  -P, --port string     port the node listens on
```

### SEE ALSO

* [stkcli](stkcli.md)	 - Manage a Statiko node

