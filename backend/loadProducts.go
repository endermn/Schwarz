package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func readProductsFromCSV(filePath string, box *productBox) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open data")
		os.Exit(1)
	}
	defer file.Close()

	err = box.RemoveAll()
	if err != nil {
		log.Printf("Failed to clear the products table: %v", err)
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV file")
		os.Exit(1)
	}

	for _, record := range records {
		id, err := strconv.Atoi(record[0][1:])
		if err != nil {
			log.Printf("strconv.Atoi failed: %v", record[0][1:])
			os.Exit(1)
		}
		_, err = box.Insert(&product{
			ProductID: id,
			Category:  record[1],
			Name:      record[2],
		})
		if err != nil {
			log.Printf("failed to insert product %v: %v", record, err)
			os.Exit(1)
		}
	}
}
