module github.com/searKing/golang/internal/tools

go 1.27.0

tool (
	github.com/searKing/golang/tools/go-atomicvalue
	github.com/searKing/golang/tools/go-enum
	github.com/searKing/golang/tools/go-import
	github.com/searKing/golang/tools/go-nulljson
	github.com/searKing/golang/tools/go-option
	github.com/searKing/golang/tools/go-sqlx
	github.com/searKing/golang/tools/go-syncmap
	github.com/searKing/golang/tools/go-syncpool
	github.com/searKing/golang/tools/go-union
	github.com/searKing/golang/tools/go-validator
	github.com/searKing/golang/tools/protoc-gen-go-tag
	go.opentelemetry.io/build-tools/crosslink
	go.opentelemetry.io/build-tools/multimod
	golang.org/x/tools/cmd/stringer
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.4.1 // indirect
	github.com/cloudflare/circl v1.6.3 // indirect
	github.com/cyphar/filepath-securejoin v0.6.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.9.0 // indirect
	github.com/go-git/go-git/v5 v5.18.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/pelletier/go-toml/v2 v2.3.0 // indirect
	github.com/pjbgf/sha1cd v0.5.0 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/searKing/golang/go v1.2.122 // indirect
	github.com/searKing/golang/tools v1.2.122 // indirect
	github.com/searKing/golang/tools/go-atomicvalue v1.2.147 // indirect
	github.com/searKing/golang/tools/go-enum v1.2.147 // indirect
	github.com/searKing/golang/tools/go-import v1.2.147 // indirect
	github.com/searKing/golang/tools/go-nulljson v1.2.147 // indirect
	github.com/searKing/golang/tools/go-option v1.2.147 // indirect
	github.com/searKing/golang/tools/go-sqlx v1.2.147 // indirect
	github.com/searKing/golang/tools/go-syncmap v1.2.147 // indirect
	github.com/searKing/golang/tools/go-syncpool v1.2.147 // indirect
	github.com/searKing/golang/tools/go-union v1.2.147 // indirect
	github.com/searKing/golang/tools/go-validator v1.2.147 // indirect
	github.com/searKing/golang/tools/protoc-gen-go-tag v1.2.147 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/skeema/knownhosts v1.3.2 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.opentelemetry.io/build-tools v0.30.0 // indirect
	go.opentelemetry.io/build-tools/crosslink v0.30.0 // indirect
	go.opentelemetry.io/build-tools/multimod v0.30.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/exp v0.0.0-20260410095643-746e56fc9e2f // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

replace github.com/searKing/golang/tools/go-atomicvalue => ../../tools/go-atomicvalue

replace github.com/searKing/golang/tools/go-enum => ../../tools/go-enum

replace github.com/searKing/golang/tools/go-import => ../../tools/go-import

replace github.com/searKing/golang/tools/go-nulljson => ../../tools/go-nulljson

replace github.com/searKing/golang/tools/go-option => ../../tools/go-option

replace github.com/searKing/golang/tools/go-sqlx => ../../tools/go-sqlx

replace github.com/searKing/golang/tools/go-syncmap => ../../tools/go-syncmap

replace github.com/searKing/golang/tools/go-syncpool => ../../tools/go-syncpool

replace github.com/searKing/golang/tools/go-union => ../../tools/go-union

replace github.com/searKing/golang/tools/go-validator => ../../tools/go-validator

replace github.com/searKing/golang/tools/protoc-gen-go-tag => ../../tools/protoc-gen-go-tag

replace github.com/searKing/golang/tools => ../../tools

replace github.com/searKing/golang/go => ../../go
