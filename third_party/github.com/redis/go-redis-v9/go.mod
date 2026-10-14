// Deprecated: Use github.com/searKing/golang/pkg/log/redislog instead.
module github.com/searKing/golang/third_party/github.com/redis/go-redis-v9

go 1.27.0

require (
	github.com/redis/go-redis/v9 v9.22.0
	github.com/searKing/golang/pkg/log/redislog v1.2.147
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/searKing/golang/pkg/log/redislog => ../../../../pkg/log/redislog
