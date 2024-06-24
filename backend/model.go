package main

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type user struct {
	id       uint64
	username string
	password string // TODO: maybe hash it?
	// isAdmin  bool
}

type store struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
