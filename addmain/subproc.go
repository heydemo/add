package addmain

import (
	"log"
	"os"
	"os/exec"
)

func Subproc(command string, args ...string) {
	cmd := exec.Command(command, args...)

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
