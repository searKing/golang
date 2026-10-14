module github.com/searKing/golang/third_party/github.com/go-sql-driver/mysql

go 1.27.0

require (
	github.com/go-sql-driver/mysql v1.10.1
	github.com/searKing/golang/go v1.2.149
)

require filippo.io/edwards25519 v1.2.0 // indirect

replace github.com/searKing/golang/go => ../../../../go

// v1.2.148 was published with incorrect commit author metadata.
retract v1.2.148
