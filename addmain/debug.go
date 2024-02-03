package addmain

import (
	"bufio"
	"fmt"
	"os"
)

func Debug() {
	if os.Getenv("ADD_DEBUG") == "true" {
		fmt.Println("Waiting for debugger to attach. Press ENTER to continue...")
		fmt.Printf("Current Process ID: %d\n", os.Getpid())
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
