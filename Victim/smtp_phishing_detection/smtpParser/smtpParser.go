package smtpparser

import (
	"bufio"
	"bytes"
	"log"
	"net/textproto"
	"strings"

	"github.com/google/gopacket/layers"
)

func SmtpBodyReader(tcp *layers.TCP) string {
	bufReader := bufio.NewReader(bytes.NewReader(tcp.BaseLayer.Payload))

	tpReader := textproto.NewReader(bufReader)

	headers, _ := tpReader.ReadMIMEHeader()
	// if err != nil {
	// 	// log.Println(err)
	// }

	// Get the subject from the headers
	_, ok := headers["Subject"]
	if ok {
		body, err := tpReader.ReadDotBytes()
		if err != nil {
			log.Println(err)
		}
		updatedBody := getMailBody(string(body))
		return updatedBody
		// fmt.Println(string(body))
	}
	return ""

}

func getMailBody(body string) string {
	lines := strings.Split(body, "\n")

	i := 0

	for !strings.HasPrefix(lines[i], "Content-Transfer-Encoding:") {
		i++
	}
	updatedLines := lines[i+2 : len(lines)-3]

	desiredBody := strings.Join(updatedLines, "\n")
	return desiredBody

}
