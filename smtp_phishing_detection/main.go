package main

import (
	"fmt"
	"smtp_phishing_detection/onlineTrafficReader"
	"smtp_phishing_detection/yara"
)

func main() {
	yara.YaraInit()
	fmt.Println("System Started.")
	onlineTrafficReader.OnlineTrafficReader()
}
