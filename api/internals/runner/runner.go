package runner

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const PIPELINE_TIMEOUT_MINUTES int = 10
const PIPELINE_ENTRY_SCRIPT_NAME string = "./start.sh"

type LastWorkflowStat struct {
	Start  time.Time `json:"last_wf_start"`
	Finish time.Time `json:"last_wf_finish"`
}

type Runner struct {
	LastWorkflowSinceStart LastWorkflowStat
	IsWorkflowRunning      bool
}

func (r *Runner) ExecutePipeline(ctx context.Context, shutdownChan *chan bool) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(PIPELINE_TIMEOUT_MINUTES)*time.Minute)
	defer cancel()

	defer func() {
		r.IsWorkflowRunning = false
	}()

	// get gracefulShutdownChan

	pwd, err := os.Getwd()
	if err != nil {
		// TODO: log error
		return fmt.Errorf("error getting working directory: %v", err)
	}

	slicedPwd := strings.Split(pwd, "/")
	poppedPwd := slicedPwd[:len(slicedPwd)-1]
	newPwd := strings.Join(poppedPwd, "/")

	cmd := exec.Command(PIPELINE_ENTRY_SCRIPT_NAME)
	cmd.Dir = newPwd

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return err
	}
	log.Printf("pipeline started with pid %d", cmd.Process.Pid)

	done := make(chan error, 1)
	go func() {
		log.Println("started done chan goroutine")
		done <- cmd.Wait()
		log.Println("ended done chan goroutine")
		close(done)
	}()

	log.Println("before select")

	select {
	case sig := <-*shutdownChan:
		if sig {
			log.Println("pipeline shutdown signal received via shutdownChan")
			if cmd.Process != nil {
				// break execution of running pipeline for now
			}

			<-done
			return nil
		}
	case <-ctx.Done():
		log.Println("pipeline timeout or shutdown signal received, killing process...")
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				return fmt.Errorf("failed to kill process: %v", err)
			}
		}

		<-done
		return ctx.Err()

	case err := <-done:
		if err != nil {
			log.Printf("pipeline finished with an error: %v", err)
		} else {
			log.Println("pipeline finished successfully")
		}

		r.LastWorkflowSinceStart.Finish = time.Now().UTC()
		return err
	}

	return nil
}
