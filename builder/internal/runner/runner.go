package runner

import (
	"fmt"
	"os/exec"
)

type Runner interface {
	Build(projectPath string, buildCommand string) error
}

func GetRunner(framework string) Runner {
	switch framework {
	case "nextjs":
		return &NextJSRunner{}
	case "nuxtjs":
		return &NuxtRunner{}
	case "nodejs":
		return &NodeRunner{}
	case "bun":
		return &BunRunner{}
	case "go":
		return &GoRunner{}
	case "php":
		return &PHPRunner{}
	case "static":
		return &StaticRunner{}
	default:
		return &StaticRunner{}
	}
}

// Next.js Runner
type NextJSRunner struct{}

func (r *NextJSRunner) Build(projectPath string, buildCommand string) error {
	if buildCommand == "" {
		buildCommand = "npm run build"
	}

	// Install dependencies
	if err := runCommand(projectPath, "npm", "install"); err != nil {
		return err
	}

	// Build
	return runCommand(projectPath, "sh", "-c", buildCommand)
}

// Nuxt Runner
type NuxtRunner struct{}

func (r *NuxtRunner) Build(projectPath string, buildCommand string) error {
	if buildCommand == "" {
		buildCommand = "npm run build"
	}

	if err := runCommand(projectPath, "npm", "install"); err != nil {
		return err
	}

	return runCommand(projectPath, "sh", "-c", buildCommand)
}

// Node.js Runner
type NodeRunner struct{}

func (r *NodeRunner) Build(projectPath string, buildCommand string) error {
	if err := runCommand(projectPath, "npm", "install"); err != nil {
		return err
	}

	if buildCommand != "" && buildCommand != "npm run build" {
		return runCommand(projectPath, "sh", "-c", buildCommand)
	}

	return nil
}

// Bun Runner
type BunRunner struct{}

func (r *BunRunner) Build(projectPath string, buildCommand string) error {
	if err := runCommand(projectPath, "bun", "install"); err != nil {
		return err
	}

	if buildCommand != "" {
		return runCommand(projectPath, "sh", "-c", buildCommand)
	}

	return nil
}

// Go Runner
type GoRunner struct{}

func (r *GoRunner) Build(projectPath string, buildCommand string) error {
	if buildCommand == "" {
		buildCommand = "go build -o main ."
	}

	return runCommand(projectPath, "sh", "-c", buildCommand)
}

// PHP Runner
type PHPRunner struct{}

func (r *PHPRunner) Build(projectPath string, buildCommand string) error {
	// Check if composer.json exists
	if err := runCommand(projectPath, "composer", "install", "--no-dev", "--optimize-autoloader"); err != nil {
		return err
	}

	if buildCommand != "" {
		return runCommand(projectPath, "sh", "-c", buildCommand)
	}

	return nil
}

// Static Site Runner
type StaticRunner struct{}

func (r *StaticRunner) Build(projectPath string, buildCommand string) error {
	if buildCommand != "" {
		return runCommand(projectPath, "sh", "-c", buildCommand)
	}
	return nil
}

// Helper function
func runCommand(dir string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}
	return nil
}

