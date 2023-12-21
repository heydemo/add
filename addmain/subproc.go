package addmain

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func Subproc(command string, args ...string) {
	var cmd *exec.Cmd = exec.Command(command, args...)

	// Attach file descriptors
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Print("Failed to start: ", err, "\n")
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Print("Command exited with error: ", err, "\n")
		return
	}

}

func SubprocWithEnv(command string, env []string, args ...string) {
	var cmd *exec.Cmd = exec.Command(command, args...)
	cmd.Env = env

	// Attach file descriptors
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Print("Failed to start: ", err, "\n")
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Print("Command exited with error: ", err, "\n")
		return
	}

}

func SubprocAndOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Create a buffer to capture the command's standard output
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()

	// Convert the captured stdout to a string
	return stdout.String(), err

}
