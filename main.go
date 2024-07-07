package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//maping unix commands to windows commands
var CommandMap = map[string]string{
	"ls":   "dir",
	"cp":   "copy",
	"rm":   "del",
	"cat":  "type",
	"echo": "echo",
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	//The input loop
	for {

		//Shell styling
		currentWorkingDir, err := os.Getwd()
		currentDir := filepath.Base(currentWorkingDir)

		currentTime := time.Now()
		// Format time
		curTime := currentTime.Format("2006-01-02 15:04:05")

		//display shell
		fmt.Printf("%s[%s]> ", currentDir, curTime)

		//read inputs
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		err = executeCommand(input)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}

func executeCommand(input string) error {
	args := strings.Fields(input)

	// Check for empty commands
	if len(args) == 0 {
		return nil
	}

	//Accept some unix equivalent commands
	if windowsCmd, ok := CommandMap[args[0]]; ok {
		args[0] = windowsCmd
	}

	// Execute the command
	cmd := exec.Command("cmd", "/c", strings.Join(args, " "))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
