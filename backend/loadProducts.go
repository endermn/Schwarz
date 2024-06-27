package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

func readRecordsFromCSV(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open data")
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV file")
		os.Exit(1)
	}

	return records

}

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func extractURL(line string) string {
	startIndex := strings.Index(line, "http")
	if startIndex == -1 {
		return ""
	}
	endIndex := strings.Index(line[startIndex:], " ")
	if endIndex == -1 {
		return line[startIndex:]
	}
	return line[startIndex : startIndex+endIndex]
}

func loadProducts(records [][]string, box *productBox) {

	err := box.RemoveAll()
	if err != nil {
		log.Printf("Failed to clear the products table: %v", err)
		os.Exit(1)
	}
	lines, err := readFile("./image-searcher/urls.txt")
	if err != nil {
		log.Printf("Did not manage to read from file: %v", err)
	}

	for i, record := range records {
		id, err := strconv.Atoi(record[0][1:])
		if err != nil {
			log.Printf("strconv.Atoi failed: %v", record[0][1:])
			os.Exit(1)
		}
		url := extractURL(lines[i])

		activeProduct := &product{
			ProductID: id,
			Category:  record[1],
			Name:      record[2],
			ImageURL:  url,
		}
		_, err = box.Put(activeProduct)
		if err != nil {
			log.Printf("failed to insert product %v: %v", record, err)
			os.Exit(1)
		}

	}
}
