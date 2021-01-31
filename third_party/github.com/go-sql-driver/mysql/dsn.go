package mysql

import (
	"net/url"

	"github.com/go-sql-driver/mysql"

	"github.com/searKing/golang/go/database/dsn"
)

// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func ParseDSN(dsn_ string) (schema string, cfg *mysql.Config, err error) {
	schema, dsnSchemaOmitted, query := dsn.Split(dsn_)
	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	cfg, err = mysql.ParseDSN(dsn.Join("", dsnSchemaOmitted, query))
	return schema, cfg, err
}

// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func GetDSN(schema string, cfg *mysql.Config) string {
	return dsn.Join(schema, cfg.FormatDSN(), url.Values{})
}
