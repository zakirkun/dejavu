package detector

import (
	"os"
	"path/filepath"
)

// Detect mendeteksi framework berdasarkan file yang ada
func Detect(projectPath string) string {
	// Check Next.js
	if fileExists(filepath.Join(projectPath, "next.config.js")) ||
		fileExists(filepath.Join(projectPath, "next.config.mjs")) {
		return "nextjs"
	}

	// Check Nuxt
	if fileExists(filepath.Join(projectPath, "nuxt.config.js")) ||
		fileExists(filepath.Join(projectPath, "nuxt.config.ts")) {
		return "nuxtjs"
	}

	// Check Go
	if fileExists(filepath.Join(projectPath, "go.mod")) {
		return "go"
	}

	// Check PHP/Laravel
	if fileExists(filepath.Join(projectPath, "composer.json")) {
		return "php"
	}

	// Check Bun
	if fileExists(filepath.Join(projectPath, "bun.lockb")) {
		return "bun"
	}

	// Check Node.js
	if fileExists(filepath.Join(projectPath, "package.json")) {
		return "nodejs"
	}

	// Check static site
	if fileExists(filepath.Join(projectPath, "index.html")) {
		return "static"
	}

	// Default to static
	return "static"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

