package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	scriptsDir string
	triggerDir string
	lockFile string
)

var mu sync.Mutex
var wg sync.WaitGroup

func main(){
	fmt.Println("start listening")
	
	// getting current dir
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// setting up paths
	scriptsDir = pwd
	triggerDir = scriptsDir + "/trigger"
	lockFile = scriptsDir + "/.script-lock"

	// creating it directory
	os.MkdirAll(triggerDir, 0755)

	// set up done and interrupt channels
	doneChan := make(chan bool, 1)
	defer close(doneChan)

	interruptCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	// goroutine for the listener
	go func(interruptCtx context.Context){
		defer wg.Done()
		for {
			fmt.Println("an iteration...")

			files, err := filepath.Glob(triggerDir + "/*.trigger")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(files) == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			
			

			mu.Lock()
			if isLocked() {
				mu.Unlock()
				time.Sleep(1 * time.Second)
				continue
			}
	
			createLock()
			triggerFile := files[0]
			
			cmd:= exec.Command("./start.sh")
	
			var outBuf bytes.Buffer
			var errBuf bytes.Buffer
	
			cmd.Stdout = &outBuf
			cmd.Stderr = &errBuf


			err = cmd.Run()
			if err != nil {
				os.WriteFile(triggerFile + ".error", errBuf.Bytes(), 0644)
			}
	
			os.WriteFile(triggerFile + ".done", outBuf.Bytes(), 0644)
	
			removeLock()
			os.Remove(triggerFile)
			mu.Unlock()

			doneChan <- true

			// waiting for channels
			select {
			case <- doneChan:
				// workflow done
			case <- interruptCtx.Done():
				// workflow interrupted
				interruptCtx.Err()
			}
		}
	}(interruptCtx)

	go func(){
		for {
			interruptFiles, err := filepath.Glob(triggerDir + "/*.interrupt")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if len(interruptFiles) == 0 {
				// has to have an xyz.interrupt file
			}

			interruptFile := interruptFiles[0]

			interruptBytes, err := os.ReadFile(interruptFile)
			if err != nil {

			}
			if strings.Contains(string(interruptBytes), "workflow_interrupt") {
				cancel()
				return
			}

			time.Sleep(300 * time.Millisecond)
		}
	}()

	wg.Wait()
	fmt.Println("end of listener")
}

func isLocked() bool {
	_, err := os.Stat(lockFile)
	return err == nil
}

func createLock() error {
	return os.WriteFile(lockFile, []byte{}, 0644)
}

func removeLock() error {
	return os.Remove(lockFile)
}