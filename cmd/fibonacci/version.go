package main

import "fmt"

var (
	version     = "v1.0"
	releaseDate = "31.01.22"
	os          = "ubuntu 20.04"
)

func printVersion() {
	fmt.Printf("Fibonacci version=%v release_date=%v os=%v", version, releaseDate, os)
}
