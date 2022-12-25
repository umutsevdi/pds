package onlineTrafficReader

import (
	"fmt"
	"log"
	smtpparser "smtp_phishing_detection/smtpParser"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func OnlineTrafficReader() {
	handle, err := pcap.OpenLive("wlp1s0", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	err = handle.SetBPFFilter("port 1025")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			body := smtpparser.SmtpBodyReader(tcp)
			fmt.Println(body)
			//TODO: send body to py
		}
	}

}
