package main

import (
	"os"

	"github.com/jenkins-x-labs/step-go-releaser/pkg"
)

// Entrypoint for the command
func main() {
	err := pkg.NewCmdGoReleaser().Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
