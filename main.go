package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

const defaultFailedCode = 1

var wg = sync.WaitGroup{}
var results = ""

func main() {
	f, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "icap", log.LstdFlags)
	filesname := []string{"hello0", "hello1", "hello2", "hello3", "hello4", "hello5"}
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
	servername := "eu.icap.glasswall-icap.com"
	port := "1344"

	var args = make([][]string, req)
	for index, element := range filesname {

		args[index] = []string{"c-icap-client", "-i", servername, "-p", port, "-f", element + ".pdf", "-s", "gw_rebuild", "-o", element + "reb" + ".pdf", "-v"}

	}
	go func(ch <-chan []string) {

		for {
			if i, ok := <-ch; ok {
				out, runerr, _ = RunCommand("time", i)
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

func RunCommand(name string, args []string) (stdout string, stderr string, exitCode int) {
	log.Println("run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	//cmd := exec.Command(name)
	cmd := exec.Command(name, args...)
	//cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	_, err := cmd.Output()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	return
}
