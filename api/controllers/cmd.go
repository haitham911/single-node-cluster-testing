package controllers

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func processor(jobInChan <-chan string, jobOutChan chan<- bool) {
	// Wait for resources to come in.
	for job := range jobInChan {
		fmt.Print(job, "in")
		t1 := time.Now()
		fmt.Println(t1.String())
		//runcmd(job)
		fmt.Println(job, "out")
		//

		//time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

		t2 := time.Now()
		fmt.Println(t2.String())

		// Fill in `jobOutChan` channel by one resource as soon as the new resource is processed.
		jobOutChan <- true
	}
}
func runcmd(j string) {
	//	cmd := exec.Command("time c-icap-client -f ./hi.pdf -i 78.159.113.46 -p 1344 -s gw_rebuild -o ./reb.pdf -v")
	//	args[index] = []string{"c-icap-client", "-i", servername, "-p", port, "-f", pathin + element + ".pdf", "-s", "gw_rebuild", "-o", pathout + "reb_" + element + ".pdf", "-v"}

	cmd := exec.Command("time", "c-icap-client", "-f", j, "-i", "54.154.157.201", "-p", "1344", "-s", "gw_rebuild", "-o", "reb_"+j, "-v")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}
