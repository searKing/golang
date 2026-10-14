module github.com/searKing/golang/third_party/github.com/redis/go-redis-v9/test

go 1.27.0

require (
	github.com/alicebob/miniredis/v2 v2.40.0
	github.com/redis/go-redis/v9 v9.22.0
	github.com/searKing/golang/third_party/github.com/redis/go-redis-v9 v1.2.146
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/searKing/golang/third_party/github.com/redis/go-redis-v9 => ../
