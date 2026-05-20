package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ExecutePipeline(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute * 10)
	defer cancel()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	slicedPwd := strings.Split(pwd, "/")
	poppedPwd := slicedPwd[:len(slicedPwd)-1]
	newPwd := strings.Join(poppedPwd, "/")
	fmt.Println(newPwd)
	cmd := exec.Command("./start.sh")
	cmd.Dir = newPwd

	var out bytes.Buffer

	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return err
	}

	// at tis time, we don't need the workflow output, the email notification handles the details
	//log.Println("\n", out.String())
	
	return nil
}