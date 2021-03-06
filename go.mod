module github.com/afritzler/garden-universe

go 1.16

replace k8s.io/client-go => k8s.io/client-go v0.20.5

require (
	github.com/gardener/gardener v1.25.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cobra v1.2.0
	github.com/spf13/viper v1.8.1
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)
