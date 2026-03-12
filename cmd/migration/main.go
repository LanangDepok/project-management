package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/LanangDepok/project-management/config"
)

func migrationsDir() string {
	dir := "database/migrations"
	for i, arg := range os.Args {
		if arg == "--dir" && i+1 < len(os.Args) {
			dir = os.Args[i+1]
		}
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Invalid migrations dir: %v", err)
	}
	return abs
}

func main() {
	config.LoadEnv()
	config.ConnectDB()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/migration [up|down|create <name>]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		run("up")
	case "down":
		run("down")
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run ./cmd/migration create <name>")
			os.Exit(1)
		}
		create(os.Args[2])
	default:
		fmt.Printf("Unknown command %q. Use: up | down | create <name>\n", os.Args[1])
		os.Exit(1)
	}
}

func create(name string) {
	dir := migrationsDir()
	files, _ := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	next := len(files) + 1

	base := fmt.Sprintf("%06d_%s", next, name)
	upFile := filepath.Join(dir, base+".up.sql")
	downFile := filepath.Join(dir, base+".down.sql")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create migrations dir: %v", err)
	}
	os.WriteFile(upFile, []byte("-- Write your UP migration here\n"), 0644)
	os.WriteFile(downFile, []byte("-- Write your DOWN migration here\n"), 0644)

	log.Printf("Created: %s", upFile)
	log.Printf("Created: %s", downFile)
}

func run(direction string) {
	dir := migrationsDir()
	log.Printf("Running migrations (%s) from: %s", direction, dir)

	pattern := filepath.Join(dir, fmt.Sprintf("*.%s.sql", direction))
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		log.Fatalf("No %s migration files found in %s", direction, dir)
	}

	sort.Strings(files)
	if direction == "down" {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	for _, file := range files {
		log.Printf("  → %s", filepath.Base(file))
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", file, err)
		}
		for _, stmt := range strings.Split(string(content), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if err := config.DB.Exec(stmt).Error; err != nil {
				log.Fatalf("Failed to execute %s: %v", filepath.Base(file), err)
			}
		}
	}
	log.Printf("Migration %s completed.", direction)
}
