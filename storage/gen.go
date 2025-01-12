//go:build ignore
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration struct holds the filename and its content
type Migration struct {
	Number  int
	Name    string
	Content string
}

// Helper function to extract the migration number from the filename
func extractMigrationNumber(filename string) (int, error) {
	// Split the filename by underscores, get the first part and convert to int
	parts := strings.Split(filename, "_")
	var num int
	_, err := fmt.Sscanf(parts[0], "%d", &num)
	return num, err
}

func main() {
	// Define command-line flags
	migrationDir := flag.String("dir", "./migrations", "Directory containing the .up.sql migration files")
	outputFile := flag.String("out", "schema.sql", "Output file for combined migrations")

	// Parse command-line flags
	flag.Parse()

	// Slice to store migrations
	var migrations []Migration

	// Read all files in the specified directory
	files, err := os.ReadDir(*migrationDir) // os.ReadDir replaces ioutil.ReadDir
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", *migrationDir, err)
		return
	}

	// Loop over the files and process only .up.sql files
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrationNumber, err := extractMigrationNumber(file.Name())
			if err != nil {
				fmt.Printf("Error extracting migration number from %s: %v\n", file.Name(), err)
				continue
			}

			// Read the content of the migration file
			content, err := os.ReadFile(filepath.Join(*migrationDir, file.Name())) // os.ReadFile replaces ioutil.ReadFile
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			// Append the migration to the slice
			migrations = append(migrations, Migration{
				Number:  migrationNumber,
				Name:    file.Name(),
				Content: string(content),
			})
		}
	}

	// Sort the migrations by their number
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Number < migrations[j].Number
	})

	// Create or open the output schema.sql file
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", *outputFile, err)
		return
	}
	defer outFile.Close()

	// Write the comment header to the output file
	_, err = outFile.WriteString(fmt.Sprintf(`-- Code generated by MigrationCombiner. DO NOT EDIT.
-- Source: %s
--
-- Generated by this command:
--
--  go run cmd/genmin/main.go -dir=%s -out=%s
--
-- Generation timestamp: %s
--
`, *migrationDir, *migrationDir, *outputFile, time.Now().Format(time.RFC1123)))
	if err != nil {
		fmt.Printf("Error writing comment header to output file %s: %v\n", *outputFile, err)
		return
	}

	// Write the content of each migration to the output file
	for _, migration := range migrations {
		fmt.Printf("Writing migration: %s\n", migration.Name)
		_, err := outFile.WriteString("-- " + migration.Name + "\n")
		if err != nil {
			fmt.Printf("Error writing to output file %s: %v\n", *outputFile, err)
			return
		}
		_, err = outFile.WriteString(migration.Content + "\n\n")
		if err != nil {
			fmt.Printf("Error writing to output file %s: %v\n", *outputFile, err)
			return
		}
	}

	fmt.Printf("Migrations combined into %s successfully!\n", *outputFile)
}
