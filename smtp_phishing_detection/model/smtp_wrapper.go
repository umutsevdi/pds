package model

// #include "include/smtp.h"
import "C"
import (
	"unsafe"

	"github.com/google/gopacket"
)

type SmtpModel struct {
	MailSender   string
	MailRecipent string
	MailDateTime string
	MailSubject  string
	MailContent  string
	SmtpModel_C  C.struct_SMTP
}

func (t SmtpModel) getSMTP(packet gopacket.Packet) C.struct_SMTP {
	return C.struct_SMTP{}
}

func (t *SmtpModel) GetData(packet gopacket.Packet, payloadLen int) bool {

	// app.Payload()[0]
	C.userData((*C.struct_SMTP)(unsafe.Pointer(&t.SmtpModel_C)), (*C.u_char)(unsafe.Pointer(&packet.ApplicationLayer().Payload()[0])), C.int(payloadLen))

	if len(C.GoString(&t.SmtpModel_C.mailSender[0])) > 1 || len(C.GoString(&t.SmtpModel_C.mailRecipent[0])) > 1 {
		return true
	}
	return false
}

func (t *SmtpModel) SmtpMapper() {
	t.MailSender = C.GoString(&t.SmtpModel_C.mailSender[0])
	t.MailRecipent = C.GoString(&t.SmtpModel_C.mailRecipent[0])
}
