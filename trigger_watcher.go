package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var (
	scriptsDir string
	triggerDir string
	lockFile   string
)

var mu sync.Mutex
var wg sync.WaitGroup

func main() {
	fmt.Println("start listening")

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scriptsDir = pwd
	triggerDir = scriptsDir + "/trigger"
	lockFile = scriptsDir + "/.script-lock"

	os.MkdirAll(triggerDir, 0755)

	osShutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	doneChan := make(chan bool, 1)
	defer close(doneChan)

	interruptChan := make(chan any, 1)
	defer close(interruptChan)

	// goroutine for the listener
	go func(interruptChan <-chan any) {
		for {
			fmt.Println("trigger watcher...")

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

			cmd := exec.Command("./start.sh")

			var outBuf bytes.Buffer
			var errBuf bytes.Buffer

			cmd.Stdout = &outBuf
			cmd.Stderr = &errBuf

			err = cmd.Run()
			if err != nil {
				os.WriteFile(triggerFile+".error", errBuf.Bytes(), 0644)
			}

			os.WriteFile(triggerFile+".done", outBuf.Bytes(), 0644)

			removeLock()
			os.Remove(triggerFile)
			mu.Unlock()

			doneChan <- true

			select {
			case <-doneChan:
				// workflow done
				fmt.Println("trw select doneChan")
			case <-interruptChan:
				// workflow interrupted
				fmt.Println("trw select interruptChan")
				return
			case <-osShutdownCtx.Done():
				fmt.Println("trw select isShutdownCtx")
				return
			}
		}
	}(interruptChan)

	// goroutine for the interrupt
	go func(interruptChan chan<- any) {
		for {

			fmt.Println("interrupt watcher...")
			interruptFiles, err := filepath.Glob(triggerDir + "/*.interrupt")
			if err != nil {
				panic(err)
			}

			if len(interruptFiles) == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			interruptChan <- true

			if err := os.RemoveAll("*.interrupt"); err != nil {
				panic(err)
			}
		}
	}(interruptChan)

	<-osShutdownCtx.Done()
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
