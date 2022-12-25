package offlineTrafficReader

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/textproto"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func OfflineReader() {

	handle, err := pcap.OpenOffline("./mailhog.pcapng")
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {

		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {

			// check if the packet is destined for the SMTP port in our case 1025
			tcp, _ := tcpLayer.(*layers.TCP)
			if tcp.DstPort == 1025 || tcp.SrcPort == 1025 {

				bufReader := bufio.NewReader(bytes.NewReader(tcp.BaseLayer.Payload))

				tpReader := textproto.NewReader(bufReader)

				headers, err := tpReader.ReadMIMEHeader()
				if err != nil {
					// log.Println(err)
				}

				// Get the subject from the headers
				_, ok := headers["Subject"]
				if ok {
					body, err := tpReader.ReadDotBytes()
					if err != nil {
						log.Println(err)
					}
					updatedBody := getMailBody(string(body))
					fmt.Println(updatedBody)
					// fmt.Println(string(body))
				}

			}
		}
	}
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

// func smtpParser(packet gopacket.Packet) {
// 	if app := packet.ApplicationLayer(); app != nil {
// 		if len(app.Payload()) > 10 {
// 			var smtp model.SmtpModel
// 			if true {
// 				smtp.GetData(packet, len(packet.ApplicationLayer().Payload()))
// 				smtp.SmtpMapper()
// 				fmt.Printf("hello %#v\n", smtp)
// 			}

// 		}
// 	}
// }
