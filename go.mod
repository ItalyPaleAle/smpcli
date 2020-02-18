module github.com/ItalyPaleAle/stkcli

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v39.1.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.8.0
	github.com/Azure/go-autorest/autorest v0.9.5 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/dsnet/compress v0.0.1
	github.com/manifoldco/promptui v0.7.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.0-00010101000000-000000000000
	github.com/spf13/viper v1.6.2
	software.sslmate.com/src/go-pkcs12 v0.0.0-20190322163127-6e380ad96778
)

replace github.com/spf13/cobra => github.com/ItalyPaleAle/cobra v0.0.6-0.20200218001531-3f49bf32ab82
