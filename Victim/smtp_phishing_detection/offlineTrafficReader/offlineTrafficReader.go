package offlineTrafficReader

import (
	smtpparser "smtp_phishing_detection/smtpParser"

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
				smtpparser.SmtpBodyReader(tcp)
			}
		}
	}
}
