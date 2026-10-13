// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

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
			return strings.ToLower(fmt.Sprint(args[i]))
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

// commandKey returns the first key of cmd, or "" if cmd has no key.
func commandKey(cmd redis.Cmder) string {
	if keys := commandKeys(cmd); len(keys) > 0 {
		return truncateString(fmt.Sprint(cmd.Args()[keys[0]]))
	}
	return ""
}

// commandStatement returns cmd with keys and subcommands only, and other arguments replaced by "?",
// such as "set k1 ?", as db.query.text of OpenTelemetry for redis, truncated by truncateArgs.
func commandStatement(cmd redis.Cmder) string {
	args := cmd.Args()
	keys := commandKeys(cmd)
	_, hasSubcommand := subcommandCommands[cmd.Name()]
	parts := make([]any, len(args))
	for i, a := range args {
		if i == 0 || (i == 1 && hasSubcommand) || slices.Contains(keys, i) {
			parts[i] = a
			continue
		}
		parts[i] = "?"
	}
	return strings.TrimSuffix(strings.TrimPrefix(fmt.Sprint(truncateArgs(parts)), "["), "]")
}

// redacted replaces arguments and replies which may carry credentials in logs,
// as [net/url.URL.Redacted] does for passwords.
const redacted = "xxxxx"

// sensitiveArgsFrom returns the index of the first argument of cmd which may carry credentials,
// such as passwords of AUTH and requirepass of CONFIG SET, or 0 if none.
func sensitiveArgsFrom(cmd redis.Cmder) int {
	switch cmd.Name() {
	case "auth", "migrate":
		return 1
	case "hello", "acl", "config":
		// HELLO protover [AUTH username password], ACL subcommand ..., CONFIG subcommand ...
		return 2
	}
	return 0
}

// commandArgs returns arguments of cmd, with those which may carry credentials redacted,
// truncated by truncateArgs.
func commandArgs(cmd redis.Cmder) []any {
	args := slices.Clone(cmd.Args())
	if from := sensitiveArgsFrom(cmd); from != 0 {
		for i := from; i < len(args); i++ {
			args[i] = redacted
		}
	}
	return truncateArgs(args)
}

// commandReply returns the reply of cmd by its Val method, redacted if cmd may carry credentials,
// such as passwords of CONFIG GET requirepass, or the error of cmd if failed, truncated by truncateValue.
func commandReply(cmd redis.Cmder) any {
	if err := cmd.Err(); err != nil {
		return err.Error()
	}
	if sensitiveArgsFrom(cmd) != 0 {
		return redacted
	}
	if val := reflect.ValueOf(cmd).MethodByName("Val"); val.IsValid() && val.Type().NumIn() == 0 && val.Type().NumOut() == 1 {
		return truncateValue(val.Call(nil)[0].Interface())
	}
	return truncateString(cmd.String())
}

// Limits of payloads logged, as SLOWLOG of redis does, which logs up to 32 arguments of a command,
// replacing the rest by "... (N more arguments)", and up to 128 bytes of an argument,
// followed by "... (N more bytes)".
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

// truncateArgs truncates args to maxLoggedElems, and each of args by truncateValue, as SLOWLOG of redis does.
func truncateArgs(args []any) []any {
	return truncateElems(args, "arguments")
}

// truncateElems truncates elems to maxLoggedElems, replacing the rest by "... (N more what)",
// and each of elems by truncateValue.
func truncateElems[E any](elems []E, what string) []any {
	n := min(len(elems), maxLoggedElems)
	if len(elems) > maxLoggedElems {
		n = maxLoggedElems - 1
	}
	out := make([]any, 0, min(len(elems), maxLoggedElems))
	for _, e := range elems[:n] {
		out = append(out, truncateValue(e))
	}
	if n < len(elems) {
		out = append(out, fmt.Sprintf("... (%d more %s)", len(elems)-n, what))
	}
	return out
}

// truncateValue truncates strings and bytes by truncateString, and slices and maps to maxLoggedElems,
// recursively, logging bytes as strings.
func truncateValue(v any) any {
	switch v := v.(type) {
	case string:
		return truncateString(v)
	case []byte:
		return truncateString(string(v))
	case nil:
		return nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		elems := make([]any, rv.Len())
		for i := range elems {
			elems[i] = rv.Index(i).Interface()
		}
		return truncateElems(elems, "elements")
	case reflect.Map:
		type entry struct {
			key string
			val any
		}
		entries := make([]entry, 0, rv.Len())
		for it := rv.MapRange(); it.Next(); {
			entries = append(entries, entry{fmt.Sprint(it.Key().Interface()), it.Value().Interface()})
		}
		slices.SortFunc(entries, func(a, b entry) int { return cmp.Compare(a.key, b.key) })
		n := min(len(entries), maxLoggedElems)
		if len(entries) > maxLoggedElems {
			n = maxLoggedElems - 1
		}
		out := make(map[string]any, n+1)
		for _, e := range entries[:n] {
			out[truncateString(e.key)] = truncateValue(e.val)
		}
		if n < len(entries) {
			out["..."] = fmt.Sprintf("(%d more entries)", len(entries)-n)
		}
		return out
	default:
		return v
	}
}
