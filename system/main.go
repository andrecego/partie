package system

import (
	"bytes"
	"fmt"
	"os/exec"
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
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
