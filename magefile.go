// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Clean files and folders used for test data or binaries
func Clean() error {
	fmt.Println("Cleaning files and folders...")
	clean := exec.Command("rm", "-rf", "dist")
	err := clean.Run()
	if err != nil {
		return err
	}

	clean = exec.Command("rm", "-f", "bin/icon")
	err = clean.Run()
	if err != nil {
		return err
	}

	return nil
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(Clean)
	mg.Deps(InstallDeps)
	fmt.Println("Building...")

	cmd := exec.Command("go", "build", "-o", "./bin/icon")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	gopath := os.Getenv("GOPATH")
	return os.Rename("./bin/icon", gopath+"/bin/identicon")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("dep", "ensure")
	return cmd.Run()
}

// Run test that creates an identicon named Identicon
func Test() error {
	mg.Deps(Build)
	fmt.Println("Running test...")

	err := os.Mkdir("dist", os.FileMode(0777))
	if err != nil {
		return err
	}

	cmd := exec.Command("./bin/icon", "-o", "dist", "Identicon")
	return cmd.Run()
}
