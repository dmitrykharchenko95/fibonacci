package main

import "fmt"

var (
	version     = "v1.1"
	releaseDate = "02.02.22"
	os          = "ubuntu 20.04"
)

func printVersion() {
	fmt.Printf("Fibonacci version=%v release_date=%v os=%v\n", version, releaseDate, os)
}
