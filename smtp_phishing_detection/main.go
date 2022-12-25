package main

import (
	"fmt"
	"smtp_phishing_detection/onlineTrafficReader"
)

func main() {
	//yara.YaraInit()
	fmt.Println("System Started.")
	onlineTrafficReader.OnlineTrafficReader()
}
