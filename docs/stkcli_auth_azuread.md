## stkcli auth azuread

Authenticate using an Azure AD account

### Synopsis

Launches a web browser to authenticate with the Azure AD application connected to the node, then stores the authentication token. This command manages the entire authentication workflow for the user, and it requires a desktop environment running on the client's machine.

The Azure AD application is defined in the node's configuration. Users must be part of the Azure AD directory and have permissions to use the app.

Once you have authenticated with Azure AD, the client obtains an OAuth token which it uses to authorize API calls with the node. Tokens have a limited lifespan, which is configurable by the admin (stkcli supports automatically refreshing tokens when possible).


```
stkcli auth azuread [flags]
```

### Options

```
  -h, --help          help for azuread
  -n, --node string   node address or IP (required)
  -P, --port string   port the node listens on
```

### SEE ALSO

* [stkcli auth](stkcli_auth.md)	 - Authenticate with a node

