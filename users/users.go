package main

type User struct {
	ID       string
	Name     string
	Password string // store only password hash
}
