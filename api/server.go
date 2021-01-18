package api

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/haitham911/single-node-cluster-testing/api/controllers"
	"github.com/joho/godotenv"
)

var wg = sync.WaitGroup{}
var results = ""

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "icap", log.LstdFlags)
	filesname := []string{"hello0", "hello1", "hello2", "hello3", "hello4", "hello5"}
	//servername := "eu.icap.glasswall-icap.com"
	servername := os.Getenv("ICAP_HOST")
	outfile, err := os.Create("./out.txt")

	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	out := ""
	runerr := ""

	ch := make(chan []string, 50)
	wg.Add(2)

	req := len(filesname)

	port := "1344"

	var args = make([][]string, req)
	for index, element := range filesname {

		args[index] = []string{"c-icap-client", "-i", servername, "-p", port, "-f", element + ".pdf", "-s", "gw_rebuild", "-o", "reb_" + element + ".pdf", "-v"}

	}
	go func(ch <-chan []string) {

		for {
			if i, ok := <-ch; ok {
				dt := time.Now()
				logger.Println("Run Command : " + dt.String())
				out, runerr, _ = controllers.RunCommand("time", i)
				dt = time.Now()
				logger.Println("Command Result : " + dt.String())
				if runerr != "" {
					logger.Println(runerr)
				} else {
					logger.Println("successful rebuild file name : " + i[6])
					results = "\n" + results + out
				}

			} else {
				break
			}
		}
		wg.Done()

	}(ch)

	go func(ch chan<- []string) {

		for index := range filesname {

			ch <- args[index]

		}

		close(ch)
		wg.Done()
	}(ch)
	wg.Wait()

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

}
