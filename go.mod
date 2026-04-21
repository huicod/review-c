module review-c

go 1.25.9

// 共享 proto：通过 git submodule 引入 huicod/reviewapis 到 ./third_party/reviewapis。
// 开发期用 replace 指向 submodule；reviewapis 打 tag 后可移除 replace 改用 `go get ...@vX.Y.Z`。
replace github.com/huicod/reviewapis => ./third_party/reviewapis

require (
	github.com/go-kratos/kratos/contrib/registry/consul/v2 v2.0.0-20260404020628-f149714c1d54
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/google/wire v0.6.0
	github.com/hashicorp/consul/api v1.34.1
	github.com/huicod/reviewapis v0.0.0-00010101000000-000000000000
	go.uber.org/automaxprocs v1.6.0
	google.golang.org/protobuf v1.36.5
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/exp v0.0.0-20260218203240-3dfff04db8fa // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250219182151-9fdb1cabc7b2 // indirect
	google.golang.org/grpc v1.70.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
