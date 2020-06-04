## stkcli auth

Authenticate with a node

### Synopsis

The auth namespace contains the commands to authenticate stkcli with a Statiko node.

The CLI supports two authentication methods:

- `psk`: pre-shared key
  A key (passphrase) used to authenticate users. The key is stored in the node's configuration file, and is transmitted by clients in the header of API calls. Clients are authenticated if the key they send matches the one in the node's configuration.
  Note that the key is not hashed nor encrypted, so using TLS to connect to nodes is strongly recommended.

- `azuread`: Azure AD account
- `auth0`: Auth0
  Clients are authenticated by passing an OAuth token to the node in the header of API calls, as obtained from an Azure AD or Auth0 application. Accounts must be added to the services' directory to be granted permission to use the app.
  This method allows for tighter control over authorized users, and relies on authorization tokens which have a shorter lifespan.

Note that your Statiko nodes might not be configured to support all authentication methods.
If you're the admin of a Statiko node, please refer to the documentation for configuring authentication methods.

Please also note that, in lieu of authorizing stkcli with one of the commands above, you can pass the value for the Authorization header in the REST calls (either the pre-shared key or an OAuth access token) using the `NODE_KEY` environmental variable, for each command (e.g. `NODE_KEY=my-psk stkcli site list`).


### Options

```
  -h, --help   help for auth
```

### SEE ALSO

* [stkcli](stkcli.md)	 - Manage a Statiko node
* [stkcli auth auth0](stkcli_auth_auth0.md)	 - Authenticate using Auth0
* [stkcli auth azuread](stkcli_auth_azuread.md)	 - Authenticate using an Azure AD account
* [stkcli auth psk](stkcli_auth_psk.md)	 - Authenticate using a pre-shared key

