package system

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func Shell(command string) error {
	fmt.Println("Running command:", command)
	return exec.Command("bash", "-c", command).Run()
}

func ShellOut(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fmt.Println("Running command:", command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	done := make(chan error)
	go func() {
		done <- cmd.Run()
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			return err, string(stdout.Bytes()), string(stderr.Bytes())
		case <-ticker.C:
			stdErrStr := string(stderr.Bytes())
			// TODO: Move this to the youtube package
			if stdErrStr != "" {
				for _, line := range strings.Split(stdErrStr, "\n") {
					if strings.HasPrefix(line, "[youtube:search+oauth2]") {
						fmt.Println(line)
					}
				}
			}
		}
	}
}
