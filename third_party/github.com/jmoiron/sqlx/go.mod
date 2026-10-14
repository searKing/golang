module github.com/searKing/golang/third_party/github.com/jmoiron/sqlx

go 1.27.0

require github.com/searKing/golang/go v1.2.148

replace github.com/searKing/golang/go => ../../../../go

// v1.2.148 was published with incorrect commit author metadata.
retract v1.2.148
