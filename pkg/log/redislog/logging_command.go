// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redislog

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// keylessCommands are commands without keys, such as commands of the server and connections,
// of which the arguments, such as subcommands and credentials, are not keys.
// It follows keylessCommands of go-redis, which routes commands of cluster clients by the first key.
var keylessCommands = map[string]struct{}{
	"acl": {}, "asking": {}, "auth": {}, "bgrewriteaof": {}, "bgsave": {}, "client": {}, "cluster": {},
	"command": {}, "config": {}, "debug": {}, "discard": {}, "echo": {}, "exec": {}, "failover": {},
	"function": {}, "hello": {}, "latency": {}, "lolwut": {}, "module": {}, "monitor": {}, "multi": {},
	"ping": {}, "psubscribe": {}, "psync": {}, "publish": {}, "pubsub": {}, "punsubscribe": {}, "quit": {},
	"readonly": {}, "readwrite": {}, "replconf": {}, "replicaof": {}, "role": {}, "save": {}, "script": {},
	"select": {}, "shutdown": {}, "slaveof": {}, "slowlog": {}, "spublish": {}, "ssubscribe": {},
	"subscribe": {}, "sunsubscribe": {}, "swapdb": {}, "sync": {}, "time": {}, "unsubscribe": {},
	"unwatch": {}, "wait": {}, "waitaof": {},
}

// subcommandCommands are commands of which the first argument is a subcommand, such as GET of CONFIG GET,
// kept by commandStatement.
var subcommandCommands = map[string]struct{}{
	"acl": {}, "client": {}, "cluster": {}, "command": {}, "config": {}, "debug": {}, "function": {},
	"latency": {}, "memory": {}, "module": {}, "object": {}, "pubsub": {}, "script": {}, "slowlog": {},
	"xgroup": {}, "xinfo": {},
}

// keyRange is the range of keys in arguments of a command, as first, last and step of COMMAND INFO,
// where last counts from the end of arguments if negative, such as -1 for the last argument.
type keyRange struct{ first, last, step int }

// keyRanges are ranges of keys of commands with more than the first argument as keys.
// Other commands with keys have the only key at the first argument.
var keyRanges = map[string]keyRange{
	"del": {1, -1, 1}, "unlink": {1, -1, 1}, "exists": {1, -1, 1}, "touch": {1, -1, 1}, "mget": {1, -1, 1},
	"watch": {1, -1, 1}, "sinter": {1, -1, 1}, "sunion": {1, -1, 1}, "sdiff": {1, -1, 1},
	"sinterstore": {1, -1, 1}, "sunionstore": {1, -1, 1}, "sdiffstore": {1, -1, 1},
	"pfcount": {1, -1, 1}, "pfmerge": {1, -1, 1},

	"rename": {1, 2, 1}, "renamenx": {1, 2, 1}, "copy": {1, 2, 1}, "smove": {1, 2, 1}, "lmove": {1, 2, 1},
	"blmove": {1, 2, 1}, "rpoplpush": {1, 2, 1}, "brpoplpush": {1, 2, 1}, "geosearchstore": {1, 2, 1},
	"zrangestore": {1, 2, 1}, "lcs": {1, 2, 1},

	// BLPOP key [key ...] timeout
	"blpop": {1, -2, 1}, "brpop": {1, -2, 1}, "bzpopmin": {1, -2, 1}, "bzpopmax": {1, -2, 1},
	// MSET key value [key value ...]
	"mset": {1, -1, 2}, "msetnx": {1, -1, 2},
	// BITOP operation destkey key [key ...]
	"bitop": {2, -1, 1},
}

// commandKeys returns indexes of keys in arguments of cmd, nil if cmd has no key.
func commandKeys(cmd redis.Cmder) []int {
	name := cmd.Name()
	if _, ok := keylessCommands[name]; ok {
		return nil
	}
	args := cmd.Args()
	arg := func(i int) string {
		if i < len(args) {
			return strings.ToLower(argString(args[i]))
		}
		return ""
	}
	// numKeys returns indexes of numkeys keys following the argument numkeys at i.
	numKeys := func(i int) []int {
		n, _ := strconv.Atoi(arg(i))
		return indexes(i+1, i+n, 1, len(args))
	}
	switch name {
	case "eval", "evalsha", "eval_ro", "evalsha_ro", "fcall", "fcall_ro":
		// EVAL script numkeys [key ...] [arg ...]
		return numKeys(2)
	case "zinter", "zunion", "zdiff", "zintercard", "sintercard", "lmpop", "zmpop":
		// ZINTER numkeys key [key ...] ...
		return numKeys(1)
	case "blmpop", "bzmpop":
		// BLMPOP timeout numkeys key [key ...] ...
		return numKeys(2)
	case "zunionstore", "zinterstore", "zdiffstore":
		// ZUNIONSTORE destination numkeys key [key ...] ...
		return append(indexes(1, 1, 1, len(args)), numKeys(2)...)
	case "memory":
		// MEMORY USAGE key
		if arg(1) == "usage" {
			return indexes(2, 2, 1, len(args))
		}
		return nil
	case "object", "xinfo", "xgroup":
		// OBJECT ENCODING key, XINFO STREAM key, XGROUP CREATE key group id
		if arg(1) == "help" {
			return nil
		}
		return indexes(2, 2, 1, len(args))
	case "xread", "xreadgroup":
		// XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]
		for i := 1; i < len(args); i++ {
			if arg(i) == "streams" {
				return indexes(i+1, i+(len(args)-i-1)/2, 1, len(args))
			}
		}
		return nil
	case "migrate":
		// MIGRATE host port <key | ""> destination-db timeout ...
		if arg(3) == "" {
			return nil
		}
		return indexes(3, 3, 1, len(args))
	}
	if r, ok := keyRanges[name]; ok {
		last := r.last
		if last < 0 {
			last += len(args)
		}
		return indexes(r.first, last, r.step, len(args))
	}
	return indexes(1, 1, 1, len(args))
}

// indexes returns indexes from first to last by step, of arguments of n.
func indexes(first, last, step, n int) []int {
	var is []int
	for i := first; i <= last && i < n; i += step {
		is = append(is, i)
	}
	return is
}

// blockingCommands are commands which block until their timeouts, if nothing to pop.
var blockingCommands = map[string]struct{}{
	"blpop": {}, "brpop": {}, "brpoplpush": {}, "blmove": {}, "blmpop": {},
	"bzpopmin": {}, "bzpopmax": {}, "bzmpop": {}, "wait": {}, "waitaof": {},
}

// isBlocking reports whether cmd may block by design, such as BLPOP, or XREAD with BLOCK,
// of which the duration is not counted as slow, as SLOWLOG of redis does not count blocking.
func isBlocking(cmd redis.Cmder) bool {
	name := cmd.Name()
	if _, ok := blockingCommands[name]; ok {
		return true
	}
	if name != "xread" && name != "xreadgroup" {
		return false
	}
	// XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]
	for _, arg := range cmd.Args()[1:] {
		switch strings.ToLower(argString(arg)) {
		case "block":
			return true
		case "streams":
			return false
		}
	}
	return false
}

// commandStatement returns cmd with keys and subcommands only, and other arguments replaced by "?",
// such as "set k1 ?", as db.query.text of OpenTelemetry for redis, formatted by commandLine.
func commandStatement(cmd redis.Cmder) string {
	keys := commandKeys(cmd)
	_, hasSubcommand := subcommandCommands[cmd.Name()]
	return commandLine(cmd.Args(), func(i int, arg any) string {
		if i == 0 || (i == 1 && hasSubcommand) || slices.Contains(keys, i) {
			return quoteArg(truncateString(argString(arg)))
		}
		return "?"
	})
}

// redacted replaces arguments which may carry credentials in logs,
// as [net/url.URL.Redacted] does for passwords.
const redacted = "xxxxx"

// secretConfigParams are parameters of CONFIG SET of which the values are credentials.
var secretConfigParams = []string{"requirepass", "masterauth", "tls-key-file-pass", "tls-client-key-file-pass"}

// secretArgs returns indexes of arguments of cmd which are credentials, nil if none,
// as secretArgs of rediscmd of go-redis does for db.statement of redisotel,
// such as the password of AUTH, keeping others, such as the username, for debugging.
// go-redis sends HELLO with AUTH on handshakes, and AUTH on rotations of credentials, through hooks.
func secretArgs(cmd redis.Cmder) []int {
	args := cmd.Args()
	is := func(i int, want string) bool { return i < len(args) && strings.EqualFold(argString(args[i]), want) }
	var secrets []int
	switch cmd.Name() {
	case "auth":
		// AUTH [username] password
		if len(args) > 1 {
			secrets = append(secrets, len(args)-1)
		}
	case "hello":
		// HELLO protover [AUTH username password] [SETNAME clientname]
		if is(2, "auth") && len(args) > 4 {
			secrets = append(secrets, 4)
		}
	case "config":
		// CONFIG SET parameter value [parameter value ...]
		if !is(1, "set") {
			break
		}
		for i := 2; i+1 < len(args); i += 2 {
			if slices.ContainsFunc(secretConfigParams, func(param string) bool { return is(i, param) }) {
				secrets = append(secrets, i+1)
			}
		}
	case "acl":
		// ACL SETUSER username [rule ...], where rules >password, <password, #hash and !hash carry credentials.
		if !is(1, "setuser") {
			break
		}
		for i := 3; i < len(args); i++ {
			if rule := argString(args[i]); len(rule) > 1 && strings.ContainsRune("><#!", rune(rule[0])) {
				secrets = append(secrets, i)
			}
		}
	case "migrate":
		// MIGRATE host port key destination-db timeout [COPY] [REPLACE] [AUTH password | AUTH2 username password]
		// [KEYS key [key ...]]
		for i := 6; i < len(args) && !is(i, "keys"); i++ {
			switch {
			case is(i, "auth") && i+1 < len(args):
				secrets = append(secrets, i+1)
			case is(i, "auth2") && i+2 < len(args):
				secrets = append(secrets, i+2)
			}
		}
	}
	return secrets
}

// commandArgs returns arguments of cmd as a line of redis-cli, such as `set k1 "hello world"`,
// with credentials redacted, formatted by commandLine.
func commandArgs(cmd redis.Cmder) string {
	secrets := secretArgs(cmd)
	return commandLine(cmd.Args(), func(i int, arg any) string {
		if slices.Contains(secrets, i) {
			return redacted
		}
		return quoteArg(truncateString(argString(arg)))
	})
}

// commandName returns the name of cmd in lower case, with its subcommand if any, such as "config get",
// as the operation name of OpenTelemetry for redis.
func commandName(cmd redis.Cmder) string {
	name := cmd.Name()
	if _, ok := subcommandCommands[name]; ok && len(cmd.Args()) > 1 {
		return name + " " + strings.ToLower(truncateString(argString(cmd.Args()[1])))
	}
	return name
}

// commandNames returns unique names of cmds by commandName, in order, truncated to maxLoggedElems.
func commandNames(cmds []redis.Cmder) []string {
	var names []string
	for _, cmd := range cmds {
		if name := commandName(cmd); !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	return truncateLines(names)
}

// commandLine returns args formatted by format, joined by spaces as a line of redis-cli,
// truncated to maxLoggedElems, with the rest replaced by "... (N more arguments)", as SLOWLOG of redis does.
func commandLine(args []any, format func(i int, arg any) string) string {
	n := len(args)
	if n > maxLoggedElems {
		n = maxLoggedElems - 1
	}
	parts := make([]string, 0, n+1)
	for i, arg := range args[:n] {
		parts = append(parts, format(i, arg))
	}
	if n < len(args) {
		parts = append(parts, fmt.Sprintf("... (%d more arguments)", len(args)-n))
	}
	return strings.Join(parts, " ")
}

// argString returns arg as sent to the redis server, as WriteArg of go-redis encodes it,
// such as 1 for true and RFC 3339 for time.Time.
func argString(arg any) string {
	if rv := reflect.ValueOf(arg); rv.Kind() == reflect.Pointer {
		_, marshaler := arg.(encoding.BinaryMarshaler)
		switch {
		case !marshaler || rv.Type().Elem() == reflect.TypeFor[time.Time]():
			// Pointers of other types, and *time.Time, are sent as their elements, or zero values if nil.
			if rv.IsNil() {
				return argString(reflect.Zero(rv.Type().Elem()).Interface())
			}
			return argString(rv.Elem().Interface())
		case rv.IsNil() && rv.Type().Elem().Implements(reflect.TypeFor[encoding.BinaryMarshaler]()):
			// MarshalBinary of value receivers panics on nil pointers, unlike of pointer receivers.
			return ""
		}
	}
	switch arg := arg.(type) {
	case nil:
		return ""
	case string:
		return arg
	case []byte:
		return string(arg)
	case bool:
		if arg {
			return "1"
		}
		return "0"
	case float32:
		return strconv.FormatFloat(float64(arg), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(arg, 'f', -1, 64)
	case time.Time:
		return arg.Format(time.RFC3339Nano)
	case time.Duration:
		return strconv.FormatInt(arg.Nanoseconds(), 10)
	case net.IP:
		return string(arg)
	case encoding.BinaryMarshaler:
		// MarshalBinary is called again here, after go-redis has called it to send arg, as the encoded bytes are
		// not exposed by go-redis, so MarshalBinary with side effects or costs is called twice if logged,
		// which is left as is, rather than leaving such arguments out of logs.
		if b, err := arg.MarshalBinary(); err == nil {
			return string(b)
		}
	}
	return fmt.Sprint(arg)
}

// quoteArg quotes s as an argument of redis-cli, if empty or containing spaces, quotes, backslashes
// or unprintable bytes, escaped as sdscatrepr of redis does, such as "hello world" and "\x00", or s otherwise.
func quoteArg(s string) string {
	if s != "" && !strings.ContainsFunc(s, func(r rune) bool {
		return r <= ' ' || r >= 0x7f || r == '"' || r == '\'' || r == '\\'
	}) {
		return s
	}
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\', '"':
			b.WriteByte('\\')
			b.WriteByte(c)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '\a':
			b.WriteString(`\a`)
		case '\b':
			b.WriteString(`\b`)
		default:
			if c < ' ' || c >= 0x7f {
				fmt.Fprintf(&b, `\x%02x`, c)
			} else {
				b.WriteByte(c)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}

// commandLines are lines of commands of a pipeline, logged as an array by JSON,
// or lines joined by "\n" by text, as fed to redis-cli by stdin.
type commandLines []string

func (l commandLines) MarshalText() ([]byte, error) { return []byte(strings.Join(l, "\n")), nil }
func (l commandLines) MarshalJSON() ([]byte, error) { return json.Marshal([]string(l)) }

// truncateLines truncates lines of commands of a pipeline to maxLoggedElems,
// with the rest replaced by "... (N more commands)", keeping each line as is.
func truncateLines(lines []string) commandLines {
	if len(lines) > maxLoggedElems {
		n := maxLoggedElems - 1
		lines = append(lines[:n:n], fmt.Sprintf("... (%d more commands)", len(lines)-n))
	}
	return lines
}

// Limits of commands logged, as SLOWLOG of redis does, which logs up to 32 arguments of a command,
// replacing the rest by "... (N more arguments)", and up to 128 bytes of an argument,
// followed by "... (N more bytes)", and up to 32 commands of a pipeline here.
const (
	maxLoggedElems = 32
	maxLoggedBytes = 128
)

// truncateString truncates s to maxLoggedBytes, as SLOWLOG of redis does.
func truncateString(s string) string {
	if len(s) <= maxLoggedBytes {
		return s
	}
	return fmt.Sprintf("%s... (%d more bytes)", s[:maxLoggedBytes], len(s)-maxLoggedBytes)
}
