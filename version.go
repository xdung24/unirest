package main

import "fmt"

var (
	GoOs           string = "undefined"
	GoArch         string = "undefined"
	GoVersion      string = "undefined"
	AppVersion     string = "undefined"
	GitHash        string = "undefined"
	GinVersion     string = "undefined"
	GorillaVersion string = "undefined"
	EchoVersion    string = "undefined"
	FiberVersion   string = "undefined"
)

func printInfo() {
	fmt.Printf("mock-servers %s (%s) \n", AppVersion, GitHash)
	fmt.Printf("%s \n", GoVersion)
	fmt.Printf("gin: %s \n", GinVersion)
	fmt.Printf("gorilla: %s \n", GorillaVersion)
	fmt.Printf("echo: %s \n", EchoVersion)
	fmt.Printf("fiber: %s \n", FiberVersion)
}
