package concurrents

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"
)

// Request is a concurrent request. It contains an http.Request. `count` specifies the number
// of request to be made. `concurrency` defines the number of concurrency i.e how many parallel
// requests will be made.
type Request struct {
	count       int64
	concurrency int
	cmd         string
	status      int64
	mutex       sync.Mutex
}

// Response contains an http.Response and an error if occured while making the request.
type Response struct {
	cmd string
	err error
}

// HttpResponse is the getter method to get the http.Response of a concurrent Request.
func (res Response) HttpResponse() {
	return
}

// Error is the getter method to get the error of a concurrent Request.
func (res Response) Error() error {
	return res.err
}

// NewRequest is the constructor for concurrent.Request. It takes a defined http.Request,
// count which specifies the number of request to be made and concurrency of the requests.
func NewRequest(cmd string, count int64, concurrency int) (req *Request) {
	return &Request{
		count:       count,
		concurrency: concurrency,
		cmd:         cmd,
		status:      0,
	}
}

// MakeSync makes the given requests in a blocking manner and returns when all the requests
// have been completed. It returns a channel of Responses correspondong to each request.
func (req *Request) MakeSync() (res chan string) {
	res = make(chan string, req.count)
	defer close(res)

	wg := sync.WaitGroup{}

	for i := 0; i < req.concurrency; i++ {
		wg.Add(1)
		go func() {
			for {
				req.mutex.Lock()
				if req.status >= req.count {
					req.mutex.Unlock()
					break
				}

				req.status++
				req.mutex.Unlock()
				t := time.Now()
				prg := "echo"

				arg1 := t.String()
				arg2 := ""
				arg3 := ""

				cmd1 := exec.Command(prg, arg1, arg2, arg3)
				stdout, err1 := cmd1.Output()

				if err1 != nil {
					fmt.Println(err1.Error())
					return
				}

				fmt.Print(string(stdout))

				cmd := exec.Command(req.cmd, t.String())
				_, err := cmd.Output()
				newRes := string(stdout)
				if err != nil {
					log.Printf("Could not get exit code for failed program: %v, %v", "cmd", err)
					newRes = "run failed"

				} else {
					newRes = "run succse"
				}
				res <- newRes
			}

			wg.Done()

			return
		}()
	}

	wg.Wait()

	return
}

// Status returns the percentage of requests completed.
func (req *Request) Status() (completed float32) {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	fmt.Println(req.status)
	return float32(req.status) / float32(req.count) * 100
}
