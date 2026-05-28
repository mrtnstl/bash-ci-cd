package runner

import (
	"bytes"
	"context"
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

func (r *Runner) ExecutePipeline(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(PIPELINE_TIMEOUT_MINUTES)*time.Minute)
	defer cancel()

	pwd, err := os.Getwd()
	if err != nil {
		// TODO: log error
		return err
	}

	slicedPwd := strings.Split(pwd, "/")
	poppedPwd := slicedPwd[:len(slicedPwd)-1]
	newPwd := strings.Join(poppedPwd, "/")

	cmd := exec.Command(PIPELINE_ENTRY_SCRIPT_NAME)
	cmd.Dir = newPwd

	var out bytes.Buffer

	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// TODO: log error
		return err
	}

	r.LastWorkflowSinceStart.Finish = time.Now().UTC()
	r.IsWorkflowRunning = false
	// at tis time, we don't need the workflow output, the email notification handles the details
	//log.Println("\n", out.String())

	return nil
}
