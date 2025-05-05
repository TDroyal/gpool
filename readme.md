### GPool

This is a simple and efficient goroutine pool. It refers to the design idea of ants and uses the object pool sync.Pool to reuse the worker. It is a memory-friendly pool that limits the number of concurrent goroutines.

### How to use it?

#### 1.download

```bash
go get -u github.com/TDroyal/gpool

```

#### 2.example
```go
package main

import (
	"fmt"

	"github.com/TDroyal/gpool"
)

func main() {
	pool, _ := gpool.NewPool(10)
	pool.Submit(func() {
		// submit your task func
		fmt.Println("task")
	})

}
```
