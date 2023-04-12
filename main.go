package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func getBaseBranch() string {
	baseBranch := "main"

	if len(os.Args) > 1 {
		baseBranch = os.Args[1]
	}

	return baseBranch
}

func getCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
			fmt.Println(err)
	}
	return strings.TrimSpace(string(output))
}

func main() {
	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatal("Directory not found")
	}

	repoPath := currentDir

	gitPath, err := exec.LookPath("git")

	if err != nil {
		log.Fatal("git command not found")
	}

	grepPath, err := exec.LookPath("grep")

	if err != nil {
		log.Fatal("grep command not found")
	}

	currentBranch := getCurrentBranch()
	baseBranch := getBaseBranch()

	c1 := exec.Command(gitPath, "--git-dir="+repoPath+"/.git", "--work-tree="+repoPath, "branch", "--merged", baseBranch)
	c2 := exec.Command(grepPath, "-v", "-e", baseBranch, "-e", currentBranch)
	c1out, err := c1.StdoutPipe()

	if err != nil {
		log.Fatal("Error creating StdoutPipe for c1:", err)
	}

	defer c1out.Close()

	c1.Start()
	c2.Stdin = c1out

	output, err := c2.Output()

	if err != nil {
		log.Fatal("cannot find branches merged")
	}

	branchesToDelete := strings.Split(strings.TrimSpace(string(output)), "\n")
	// log.Println(branchesToDelete)

	for _, branch := range branchesToDelete {
		cmd := exec.Command(gitPath, "--git-dir="+repoPath+"/.git", "--work-tree="+repoPath, "branch", "-d", strings.TrimSpace(branch))
		err := cmd.Run()

		if err != nil {
			fmt.Printf("Error deleting branch %s: %s\n", branch, err)
		} else {
			fmt.Printf("Branch %s deleted\n", branch)
		}
	}
}
