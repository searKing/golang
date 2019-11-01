[![GoDoc](https://godoc.org/github.com/searKing/golang/go/go/generator?status.svg)](https://godoc.org/github.com/searKing/golang/go/go/generator)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/go/go/generator)](https://goreportcard.com/report/github.com/searKing/golang/go/go/generator) 
# Generator
Generator behaves like Generator in [python](https://wiki.python.org/moin/Generators) or [ES6](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*), with yield and next statements.

Generator function contains one or more yield statement.  
Generator functions allow you to declare a function that behaves like an iterator, i.e. it can be used in a for loop.  
Generator generators are a simple way of creating iterators. All the overhead we mentioned above are automatically handled by generators.  
Simply speaking, a generator is a function that returns an object (iterator) which we can iterate over (one value at a time).  
If desired, go to [python](https://wiki.python.org/moin/Generators) or [ES6](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*) for more information.  
## Example

### Golang Generators
```go
package main

import (
    "fmt"

    "github.com/searKing/golang/go/go/generator"
)

func main() {
    g := func(i int) *generator.Generator {
        return generator.GeneratorFunc(func(yield generator.Yield) {
            yield(i)
            yield(i + 10)
        })
    }

    gen := g(10)
    
    // WAY 1, by channel
    for msg := range gen.C {
        fmt.Println(msg)
    }
    // WAY 2, by Next
    //for {
    //	msg, ok := gen.Next()
    //	if !ok {
    //	    return
    //	}
    //	fmt.Println(msg)
    //}

    // Output:
    // 10
    // 20
}
```

### Python Generators
```python
# A simple generator function
def generator(i):
    n = 1
    # Generator function contains yield statements
    yield i
    yield i+10

# Using for loop
for item in generator(10):
    print(item) 

# Output:
# 10
# 20
```
### JavaScript Generators
```javascript
function* generator(i) {
  yield i;
  yield i + 10;
}

const gen = generator(10);

console.log(gen.next().value);
// expected output: 10

console.log(gen.next().value);
// expected output: 20
```

## Download/Install

The easiest way to install is to run `go get -u github.com/searKing/golang/go/go/generator`.   
You can also manually git clone the repository to `$GOPATH/src/github.com/searKing/golang/go/go/generator`.

## Inspiring Generators
* [Python](https://wiki.python.org/moin/Generators)
* [JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/function*)
