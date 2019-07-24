# dispatch

A powerful read and handle workflow dispatcher for golang.

* [Install](#install)
* [Examples](#examples)

---

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/searKing/dispatch
```

---

## Examples

Let's start registering a couple of URL paths and handlers:

```go
func main() {
    var conn chan int
    workflow := dispatch.NewDispatcher(
    	dispatch.ReaderFunc(func() (interface{}, error) {
    		return ReadMessage(conn)
    	}), 
    	dispatch.HandlerFunc(func(msg interface{}) error {
    		m := msg.(*int)
    		return HandleMessage(m)
    	}))
    workflow.Start()
}
```

Here we can set the workflow joinable. 
```go
    workflow := dispatch.NewDispatcher(nil, nil).Joinable()
    go workflow.Start()
    workflow.Join()
```

Here we can cancel the workflow. 

```go
    workflow := dispatch.NewDispatcher(nil, nil).Joinable()
    go workflow.Start()
	workflow.Context().Done()
	workflow.Join()
```

And this is all you need to know about the basic usage. More advanced options are explained below.

SEE [example](https://github.com/searKing/dispatch/tree/master/example_test.go)

---

## License

MIT licensed. See the LICENSE file for details.