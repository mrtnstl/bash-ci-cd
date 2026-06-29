package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

const PIPELINE_TIMEOUT_MINUTES int = 10
const PIPELINE_ENTRY_SCRIPT_NAME string = "start.sh"

type LastWorkflowStat struct {
	Start  time.Time `json:"last_wf_start"`
	Finish time.Time `json:"last_wf_finish"`
}

type Runner struct {
	LastWorkflowSinceStart LastWorkflowStat
	IsWorkflowRunning      bool
	PipelineScriptsLocation string
}

func (r *Runner) ExecutePipeline(ctx context.Context, shutdownChan *chan bool) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(PIPELINE_TIMEOUT_MINUTES)*time.Minute)
	defer cancel()

	defer func() {
		r.IsWorkflowRunning = false
	}()

	fmt.Println("ExecutePipeline called")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	triggerPath := r.PipelineScriptsLocation + "trigger/pipeline.trigger"
	donePath := triggerPath + ".done"
	errorPath := triggerPath + ".error"
	interruptPath := triggerPath + ".interrupt"

	if err := os.WriteFile(triggerPath, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to write trigger file: %w", err)
	}

	for {
		select {
		case <- ticker.C:
			var errs []error

			if _, err := os.Stat(donePath); err == nil {
				if _, err = os.ReadFile(donePath); err != nil {
					errs = append(errs, fmt.Errorf("failed to read done file: %w", err))
				}
				os.Remove(donePath)
			}

			if _, err := os.Stat(errorPath); err == nil {
				if _, err := os.ReadFile(errorPath); err != nil {
					errs = append(errs, fmt.Errorf("failed to read error file: %w", err))
				}
				os.Remove(errorPath)
			}
			
			return errors.Join(errs...)
		case <- ctx.Done():
			if err := os.WriteFile(interruptPath, []byte{}, 0644); err != nil {
				return fmt.Errorf("failed to write interrupt file: %w", err)
			}
			return ctx.Err()
		case <- *shutdownChan:
			cancel()
		}
	}
}
