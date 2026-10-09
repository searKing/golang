package dsn

import (
	"fmt"
	"net/url"
	"strings"
)

// See https://en.wikipedia.org/wiki/Data_source_name
// <driver>://<username>:<password>@<host>:<port>/<database>[?<param1>=<value1>&<paramN>=<valueN>]
// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]

// GetDriverName returns the driver name of a given DSN.
func GetDriverName(dsn string) string {
	scheme, _, _ := Split(dsn)
	return scheme
}

// Split splits dsn into a driver name(scheme) and left component.
// If there is no :// in dsn, Split returns an empty driver name and
// dsn without schema.
// The returned values have the property that dsn = driver+dsnSchemaOmitted.
// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
// =>
// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func Split(dsn string) (scheme string, connect string, query url.Values) {
	{
		parts := strings.SplitN(dsn, "://", 2)
		if len(parts) == 0 {
			scheme = ""
			connect = dsn
		} else if len(parts) == 1 {
			scheme = ""
			connect = parts[0]
		} else {
			scheme = parts[0]
			connect = parts[1]
		}
	}
	if i := queryIndex(connect); i >= 0 {
		query, _ = url.ParseQuery(connect[i+1:])
		connect = connect[:i]
	}
	return
}

// queryIndex returns the index of '?' starting the query in connect, or -1 if there is none.
// '?' before '@' is taken as part of password if no '/' is before it,
// like "user:pa?ss@tcp(host)/dbname", as '/' always precedes the query after host.
func queryIndex(connect string) int {
	q := strings.Index(connect, "?")
	if q < 0 || strings.Contains(connect[:q], "/") {
		return q
	}
	rest := connect[q:]
	if slash := strings.Index(rest, "/"); slash >= 0 {
		rest = rest[:slash]
	}
	at := strings.LastIndex(rest, "@")
	if at < 0 {
		return q
	}
	if i := strings.Index(connect[q+at:], "?"); i >= 0 {
		return q + at + i
	}
	return -1
}

// Join joins driver name(scheme) and dsn without schema into a single path,
// separating them with slashes.
func Join(scheme string, connect string, query url.Values) (dsn string) {
	rawQuery := query.Encode()
	if rawQuery != "" {
		connect = fmt.Sprintf("%s?%s", connect, rawQuery)
	}
	if scheme == "" {
		return connect
	}
	return fmt.Sprintf("%s://%s", scheme, connect)

}

// Masking hiding original username and password with character '*'
// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func Masking(dsn string) string {
	scheme, connect, query := Split(dsn)
	if connect == "" {
		return Join(scheme, "", query)
	}
	pos := strings.LastIndex(connect, "@")
	if pos < 0 { // No Auth
		return dsn
	}

	return Join(scheme, fmt.Sprintf("*:*@%s", connect[pos+1:]), query)
}
