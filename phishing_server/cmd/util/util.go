package util

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Configuration struct {
	lastExecutionTime int
	port              int
	token             string
}

var Config Configuration

// Reads the app.config and initializes the Configuration struct
func init() {
	file, err := ioutil.ReadFile("../app.config")
	if err != nil {
		log.Fatal("error: app.config is missing")
	}
	lines := strings.Split(string(file), "\n")
	for i := range lines {
		line := strings.Split(lines[i], "=")
		if len(line) <= 1 {
			continue
		}
		switch line[0] {
		case "port":
			Config.port, err = strconv.Atoi(line[1])
			if err != nil {
				log.Fatal("error: port must be an integer")
			}
		}
	}
	if Config.port == 0 {
		log.Println("Port parameter was not found in config file, continuing with the :8080")
		Config.port = 8080
	}
	log.Println("app.config was parsed successfuly")
}

// @{link Configuration} related Functions

// Returns the port address
// @return port to run
func (c *Configuration) Port() int {
	return c.port
}

// Returns Github Authentication Token
// @return String that has been registered
func (c *Configuration) Token() string {
	return c.token
}

// Returns the requested file in the /public/ directory
// @param url string path to file
// @return []byte content of the file, if file exist, nil otherwise
// @return nil if the file exists, error type otherwise
func ReadByteFrom(url string) ([]byte, error) {
	file, err := ioutil.ReadFile("../" + url)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Parses received url and extracts it's extension
// @param url path to file
// @return string corresponding extension type in the format of .type
func Ext(url string) string {
	p := strings.Split(url, "/")
	fname := strings.Split(p[len(p)-1], ".")
	return "." + fname[len(fname)-1]
}
