package main

import (
	"fmt"
	"smtp_phishing_detection/offlineTrafficReader"
)

func main() {
	fmt.Println("System Started.")
	offlineTrafficReader.OfflineReader()
}
