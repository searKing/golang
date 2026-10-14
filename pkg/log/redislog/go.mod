module github.com/searKing/golang/pkg/log/redislog

go 1.27.0

require github.com/redis/go-redis/v9 v9.23.0

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
)

// v1.2.148 was published with incorrect commit author metadata.
retract v1.2.148
