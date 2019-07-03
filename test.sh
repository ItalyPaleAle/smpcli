go run . -k -n localhost auth
# Shared Key: hello world

go run . -k -n localhost status

go run . -k -n localhost site add \
    -d site1.local \
    -a site1-alias.local \
    -a mysite.local \
    -c site1