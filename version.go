package main

import "fmt"

var (
	AppVersion string = "undefined"
	GitHash    string = "undefined"
	BuildTime  string = "undefined"
)

func printInfo() {
	fmt.Printf("unirest %s (%s) \n", AppVersion, GitHash)
	fmt.Printf("build time: %s \n", BuildTime)
}
