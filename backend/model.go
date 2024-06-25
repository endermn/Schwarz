package main

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

type user struct {
	id       uint64
	username string
	password string // TODO: maybe hash it?
	// isAdmin  bool
}

type product struct {
	id       uint64
	category string
	name     string
}
