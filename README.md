# signal

The signal package provides channel based synchronization primitives


###### Install
```sh
go get github.com/LydinaPavuh/signal
```
## Example

```go
package main

import (
	"fmt"
	"github.com/LydinaPavuh/signal"
	"time"
)

func backgroundWorker(mustEnd *signal.Flag) {
	waiter := mustEnd.Subscribe()
	defer waiter.Cancel()
	for {
		select {
		case <-time.Tick(1 * time.Second):
			fmt.Println("Tick from worker")
		case <-waiter.Wait():
			fmt.Println("Worker stopped")
		}
	}
}

func main() {
	mustEndFlag := signal.NewFlag()
	go backgroundWorker(mustEndFlag)
	time.Sleep(3 * time.Second)
	mustEndFlag.Raise()
	time.Sleep(1 * time.Second)
}
```
```
>>>
Tick from worker
Tick from worker
Worker stopped
```