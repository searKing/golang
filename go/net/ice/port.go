package ice

import "fmt"

var portMap = map[string]string{
	"stun":  "3478",
	"turn":  "3478",
	"stuns": "5349",
	"turns": "5349",
}
var getDefaultPort = func(schema string) (string, error) {
	port, ok := portMap[schema]
	if ok {
		return port, nil
	}
	return "", fmt.Errorf("malformed schema:%s", schema)
}
