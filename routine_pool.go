//A simple goroutine pool implemention

package gpool

import (
	"sync"
	"reflect"
	"errors"
)


//Result for the execution task
type ExecResult struct{
	Err error
	Result []interface{}
}


var initOnce sync.Once

var ErrWrongFuncType = errors.New("Wrong type for Func-param")


//struct for communication between gpool and worker function
type reqMsg struct {
	RspChan chan ExecResult  //channel for returning result
	Func interface{}
	Args []interface{}
}

type RoutinePool struct {
	poolSize uint32
	outChan chan *reqMsg  //channel for sending msg to worker function
}


//init funciton
//should be called before used
//@param size : the size of goroutine inside the pool,size is fixed after initialized
func (routinePool * RoutinePool) Init(size uint32)  {
	initOnce.Do(func() {
		if size == 0 {
			//at least one goroutine in the pool
			size = 1
		}

		routinePool.poolSize = size
		routinePool.outChan = make(chan *reqMsg)

		for i := uint32(0); i < size; i++ {
			go workerFunc(routinePool.outChan)
		}

	})
}

//TODO timeout returning

//Execute the task with response parametres
//param：@fnt executing function
//	@rspChan channel for returning results, should be buffered
//	@args the arguments  for the function
func (routinePool * RoutinePool) ExecWithRespond(fnt interface{},
				rspChan chan ExecResult, args ... interface{})  {
	var req reqMsg
	req.Func = fnt
	req.Args = args
	req.RspChan = rspChan
	routinePool.outChan <- &req
}

//Execute the task with no response parametres
//param：@fnt executing function
//	@args the arguments  for the function
func (routinePool * RoutinePool) ExecWithoutRespond(fnt interface{},args ... interface{})  {
	var req reqMsg
	req.Func = fnt
	req.Args = args
	req.RspChan = nil
	routinePool.outChan <- &req
}


//The worker function
func workerFunc(inChan chan *reqMsg)  {
	for {
		reqPtr := <-inChan
		req := *reqPtr

		funcValue := reflect.ValueOf(req.Func)
		if funcValue.Kind() != reflect.Func  {
			respondError(reqPtr,ErrWrongFuncType)
		}

		argValue := make([]reflect.Value,0)
		for i := 0; i < len(req.Args); i++ {
			argValue = append(argValue,reflect.ValueOf(req.Args[i]))
		}

		res := funcValue.Call(argValue)

		var rspErr error
		var rspResult []interface{}
		if l := len(res); l > 0{
			if err, ok := res[l - 1].Interface().(error); ok {
				rspErr = err
				for i := 0; i < l - 1; i++ {
					rspResult = append(rspResult, res[i].Interface())
				}
			}else {
				for i := 0; i < l ; i++ {
					rspResult = append(rspResult, res[i].Interface())
				}
			}

		}
		respondResult(reqPtr,rspResult,rspErr)

	}
}


func respondError(req *reqMsg, err error)  {
	if req.RspChan != nil{
		var rsp ExecResult
		rsp.Err = err
		rsp.Result = nil

		req.RspChan <- rsp
	}

}

func respondResult(req *reqMsg, results []interface{}, err error)  {
	if req.RspChan != nil {
		var rsp ExecResult
		rsp.Result = results
		rsp.Err = err

		req.RspChan <- rsp
	}

}
