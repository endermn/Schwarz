package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func searchImage(query string) string {
	searchURL := fmt.Sprintf("https://www.bing.com/images/search?q=%s", url.QueryEscape(query))

	// Make a request to Bing
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatalf("Failed to get search results: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse search results: %v", err)
	}

	// Find the first image result
	var imageURL string
	doc.Find(".mimg").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if src, exists := s.Attr("src"); exists {
			imageURL = src
			return false // Break the loop after finding the first image
		}
		return true
	})

	// Ensure the imageURL is not empty
	if imageURL == "" {
		log.Println("No images found.")
	}
	return imageURL
}

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
		url := searchImage(record[2])
		_, err = box.Insert(&product{
			ProductID: id,
			Category:  record[1],
			Name:      record[2],
			Image:     url,
		})
		if err != nil {
			log.Printf("failed to insert product %v: %v", record, err)
			os.Exit(1)
		}
	}
}
