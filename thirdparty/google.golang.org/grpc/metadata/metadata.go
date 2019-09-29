package metadata

import (
	"google.golang.org/grpc/metadata"
	"strings"
)

func New(k string, vals ...string) metadata.MD {
	md := metadata.MD{}
	key := strings.ToLower(k)
	for _, val := range vals {
		md[key] = append(md[key], val)
	}
	return md
}
