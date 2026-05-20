package cmd

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ExecutePipeline(ctx context.Context, app *Application) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute * 10)
	defer cancel()

	pwd, err := os.Getwd()
	if err != nil {
		// TODO: log error
		return err
	}

	slicedPwd := strings.Split(pwd, "/")
	poppedPwd := slicedPwd[:len(slicedPwd)-1]
	newPwd := strings.Join(poppedPwd, "/")
	
	cmd := exec.Command("./start.sh")
	cmd.Dir = newPwd

	var out bytes.Buffer

	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// TODO: log error
		return err
	}

	app.LastWorkflowSinceStart.Finish = time.Now().UTC()
	app.IsWorkflowRunning = false
	// at tis time, we don't need the workflow output, the email notification handles the details
	//log.Println("\n", out.String())
	
	return nil
}