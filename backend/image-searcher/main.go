package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func readProductsFromCSV(filePath string) [][]string {
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

func searchImage(query string, width, height int) string {
	searchURL := fmt.Sprintf("https://www.bing.com/images/search?q=%s&qpvt=%s&FORM=IRFLTR&first=1&tsc=ImageBasicHover&cw=%d&ch=%d",
		url.QueryEscape(query), url.QueryEscape(query), width, height)

	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatalf("Failed to get search results: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse search results: %v", err)
	}

	var imageURL string
	doc.Find(".mimg").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if src, exists := s.Attr("src"); exists {
			imageURL = src
			return false
		}
		return true
	})

	if imageURL == "" {
		log.Println("No images found for query:", query)
	} else {
		log.Println("Found image URL for query:", query, "->", imageURL)
	}
	return imageURL
}

func main() {
	records := readProductsFromCSV("product_master_data.csv")
	f, err := os.Create("image-searcher/urls.txt")
	if err != nil {
		panic(err)
	}
	results := make([]string, len(records))
	var wg sync.WaitGroup
	wg.Add(len(records))
	for i, record := range records {
		go func() {
			url := searchImage(record[2]+" ("+record[1]+")", 500, 500)
			results[i] = record[2] + " : " + url + "\n"
			wg.Done()
		}()
	}
	wg.Wait()

	for _, l := range results {
		l, err := f.WriteString(l)
		if err != nil {
			log.Println(err)
			f.Close()
			return
		}
		log.Println(l, "bytes written successfully")
	}
	err = f.Close()
	if err != nil {
		log.Printf("Failed to close file: %v", err)
	}
}
