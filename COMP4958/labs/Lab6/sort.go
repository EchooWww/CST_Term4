package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Record struct to store the record information
type Record struct {
	FirstName string
	LastName  string
	Score     int
}

func main() {
	// Check command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run sort.go <filename>")
		os.Exit(1)
	}

	// Open the specified file
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	// Read records from the file
	var records []Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		record, valid := parseRecord(line)
		if valid {
			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Sort the records
	sortRecords(records)

	// Print the sorted records
	for _, record := range records {
		fmt.Printf("%d %s %s\n", record.Score, record.LastName, record.FirstName)
	}
}

// parseRecord parses a line and returns a Record if valid
func parseRecord(line string) (Record, bool) {
	// Split the line into words
	words := strings.Fields(line)
	
	// Check if there are at least 3 words
	if len(words) < 3 {
		return Record{}, false
	}

	// Try to parse the third word as an integer
	score, err := strconv.Atoi(words[2])
	if err != nil || score < 0 || score > 100 {
		return Record{}, false
	}

	// Create and return a valid record
	return Record{
		FirstName: words[0],
		LastName:  words[1],
		Score:     score,
	}, true
}

// sortRecords sorts the records according to the specified rules
func sortRecords(records []Record) {
	sort.Slice(records, func(i, j int) bool {
		// Primary sort: descending order of scores
		if records[i].Score != records[j].Score {
			return records[i].Score > records[j].Score
		}
		
		// Secondary sort: ascending order of last names
		if records[i].LastName != records[j].LastName {
			return records[i].LastName < records[j].LastName
		}
		
		// Tertiary sort: ascending order of first names
		return records[i].FirstName < records[j].FirstName
	})
}