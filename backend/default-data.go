package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
<<<<<<< HEAD
	"slices"
=======
>>>>>>> main
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User

func makeDefaultUser(userBox *userBox) *user {
	_, err := userBox.Query(user_.username.Equals("admin", true)).Remove()
	if err != nil {
		panic(err)
	}
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	defaultUser := &user{username: "admin", passwordHash: passwordHash}
	_, err = userBox.Insert(defaultUser)
	if err != nil {
		panic(err)
	}
	return defaultUser
}

// Products

func readCSV(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open %v: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("failed to read %v: %v", path, err)
	}

	return records
}

func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
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

func loadProducts(productBox *productBox, path string, defaultUser *user) {
	records := readCSV(path)

	err := productBox.RemoveAll()
	if err != nil {
		log.Fatalf("failed to clear the products table: %v", err)
	}
	lines, err := readFileLines("./image-searcher/urls.txt")
	if err != nil {
		log.Printf("failed to read from file: %v", err)
	}

	for i, record := range records {
		id, err := strconv.Atoi(record[0][1:])
		if err != nil {
			log.Fatalf("strconv.Atoi failed: %v", record[0][1:])
		}
		url := extractURL(lines[i])

		activeProduct := &product{
			ProductID:   id,
			Category:    record[1],
			Name:        record[2],
			ImageURL:    url,
			IsGoldenEgg: slices.Contains(eggs, id),
			Owner:       defaultUser,
		}
		_, err = productBox.Put(activeProduct)
		if err != nil {
			log.Fatalf("failed to insert product %v: %v", record, err)
		}
	}
}

// Store

func loadDefaultStore(storeBox *storeBox, csvPath string) uint64 {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("failed to open %v: %v", csvPath, err)
	}
	defer file.Close()

	grid, start, err := parseStoreCSV(file)
	if err != nil {
		log.Fatalf("failed to parse store CSV: %v", err)
	}
	_, err = storeBox.Query(store_.Name.Equals("default", true)).Remove()
	if err != nil {
		panic(err)
	}
	defaultStoreID, err := storeBox.Insert(&store{
		Name:  "default",
		width: getWidth(grid),
		grid:  encodeGrid(grid),
		start: start,
		owner: 0,
	})
	if err != nil {
		log.Fatalf("failed to insert store into database: %v", err)
	}

	return defaultStoreID
}
