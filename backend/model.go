package main

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type user struct {
	id       uint64
	username string
	password string // TODO: maybe hash it?
	// isAdmin  bool
}

type product struct {
	id       uint64 `json:"id`
	name     string `json:"name"`
	price    int    `json:"price"`
	category string `jsoon:"category"`
}

type store struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
