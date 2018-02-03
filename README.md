# gpool
A simple goroutine pool implementation


### Init

```go
var pool gpool.RoutinePool
pool.Init(3)
``` 
Initialize the pool, the  `size`  parameter is the number of the goroutines .
The initialization will only be done for once, no matter how many times `pool.Init(num)` is called.

### ExecWithRespond & ExecWithoutRespond
When the executing function has out parameter, use the `ExecWithRespond`, otherwise the `ExecWithoutRespond`. When using the `ExecWithRespond	`, the second parameter is the channel from which we get the result.
<br/>
```go 
package main
import "fmt:"
import "gpool"

func Add(a, b int) {
	fmt.Println("Inside the add")
	return a + b
}

func main() {
	rspChan := make(chan gpool.ExecResult,3)
	pool.ExecWithoutRespond(func() {
		fmt.Println("Inside the task")
		//Do something
	})
	pool.ExecWithRespond(Add,rspChan, 1, 2)
	pool.ExecWithRespond(Add, rspChan, 4, 6)
	pool.ExecWithRespond(Add, rspChan, 5, 9)

	//get result from the channel
	for i := 0; i < 3; i ++ {
		result  := <-rspChan
		if result.Err != nil {
			fmt.Println("Error occurs, err is ", result.Err)
		}else {
			if intResult, ok := result.Result[0].(int); ok {
				fmt.Println("result is ", intResult)
			}
		}
	}
}
``` 


### ExecResult
The struct for the result. 

```go
type ExecResult struct {
	Err error
	Result [] interface{}
}
``` 



