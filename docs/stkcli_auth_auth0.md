## stkcli auth auth0

Authenticate using Auth0

### Synopsis

Launches a web browser to authenticate with the Auth0 application connected to the node, then stores the authentication token. This command manages the entire authentication workflow for the user, and it requires a desktop environment running on the client's machine.

The Auth0 application is defined in the node's configuration. Users must be part of the Auth0 directory and have permissions to use the app.

Once you have authenticated with Auth0, the client obtains an OAuth token which it uses to authorize API calls with the node. Tokens have a limited lifespan, which is configurable by the admin (stkcli supports automatically refreshing tokens when possible).


```
stkcli auth auth0 [flags]
```

### Options

```
  -h, --help          help for auth0
  -N, --node string   node address or IP
  -P, --port string   port the node listens on
```

### SEE ALSO

* [stkcli auth](stkcli_auth.md)	 - Authenticate with a node

