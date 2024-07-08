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
	"ls":  "dir",
	"cp":  "copy",
	"rm":  "del",
	"cat": "type",
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
		branch, bErr := getCurrentGitBranch()

		if bErr != nil {
			fmt.Printf("%s[%s]> ", currentDir, curTime)
		} else {
			fmt.Printf("%s (%s) [%s]> ", currentDir, branch, curTime)
		}

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

func findGitRepo() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break // Reached the root directory
		}
		dir = parentDir
	}

	return "", fmt.Errorf("no git repository found")
}

func getCurrentGitBranch() (string, error) {

	repoDir, err := findGitRepo()

	if err != nil {
		return "", err
	}
	cmd := exec.Command("git", "-C", repoDir, "rev-parse", "--abbrev-ref", "HEAD")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running git command: %v", err)
	}

	branch := strings.TrimSpace(string(output))

	return branch, nil
}
