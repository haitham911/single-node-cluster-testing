package api

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/single-node-cluster-testing/api/controllers"

	"github.com/single-node-cluster-testing/api/concurrents"

	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

var results = ""

func Path() (path string) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	fmt.Println("BEGIN")
	filesname := []string{"hello0", "hello1", "hello2", "hello3", "hello4", "hello5"}
	pathin := Path() + "/inputs/"
	pathout := Path() + "/output/"

	//servername := "eu.icap.glasswall-icap.com"  54.154.157.201
	servername := os.Getenv("ICAP_HOST")

	ch := make(chan bool, len(filesname))
	//wg.Add(2)
	wg.Add(len(filesname))

	port := "1344"

	for _, filename := range filesname {

		runcmd := []string{"c-icap-client", "-i", servername, "-p", port, "-f", pathin + filename + ".pdf", "-s", "gw_rebuild", "-o", pathout + "reb_" + filename + ".pdf", "-v"}

		go process(filename, ch, &wg, runcmd)
	}

	wg.Wait()
	fmt.Println("END")

}
func process(filename string, ch chan bool, wg *sync.WaitGroup, runcmd []string) {

	// As soon as the current goroutine finishes (job done!), notify back WaitGroup.
	defer wg.Done()

	// Acquire 1 resource to fill in the channel buffer. Once the channel buffer is full, it blocks `range` loop.
	ch <- true
	outfile, err := os.Create("./out.txt")

	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "icap", log.LstdFlags)
	dt := time.Now()
	fmt.Println("Run Command : " + filename + "  Time : " + dt.String())

	//time.Sleep(time.Duration(5) * time.Second)

	logger.Println("Run Command : " + filename)
	runerr := ""
	out := ""
	out, runerr, _ = controllers.RunCommand("time", runcmd)
	dt = time.Now()
	logger.Println("Command Result : " + dt.String())
	if runerr != "" {
		logger.Println(runerr)
	} else {
		logger.Println("successful rebuild file name : " + filename)
		results = "\n" + results + out
	}

	_, err2 := outfile.WriteString(results)
	if err2 != nil {
		panic(err2)
	}
	_, err3 := outfile.WriteString(runerr)
	if err3 != nil {
		panic(err3)
	}
	outfile.Sync()
	w := bufio.NewWriter(outfile)
	w.Flush()

	<-ch
}

func Test() {
	cmd := "echo"

	// Parallelism of the request
	concurrency := 20

	// Total number of requests to be made.
	numberOfRequests := int64(100)
	concurrentRequest := concurrents.NewRequest(cmd, numberOfRequests, concurrency)

	startTime := time.Now()
	go func() {
		concurrentRequest.MakeSync()
		completetionTime := time.Now().Sub(startTime)
		fmt.Printf("%v time required to complete all requests", completetionTime)
	}()

	tick := time.NewTicker(500 * time.Millisecond)
	for range tick.C {
		status := concurrentRequest.Status()
		timeElapsed := time.Now().Sub(startTime)
		fmt.Printf("%f% requests sent, Time elapsed: %v", status, timeElapsed)
		time.Sleep(1 * time.Second)

		// Calling Stop() method
		tick.Stop()

	}

}
