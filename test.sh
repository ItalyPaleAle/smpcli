go run . \
    auth \
    -k -n localhost
# Shared Key: hello world

go run . \
    status \
    -k -n localhost

go run . \
    site add \
    -k -n localhost \
    -d site1.local \
    -a site1-alias.local \
    -a mysite.local \
    -c site1

go run . \
    site list \
    -k -n localhost

go run . \
    site get \
    -k -n localhost \
    -d mysite.local

go run . \
    config set \
    AzureStorageContainer \
    "smplatformcistorage"

go run . \
    config get \
    AzureStorageContainer
