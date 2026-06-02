package smtp

import (
	"fmt"
	"net"
	"net/smtp"
)

type SmtpSender struct {
	host string
	port string
	from string
}

func NewSmtpSender(host string, port string, from string) *SmtpSender {
	return &SmtpSender{host: host, port: port, from: from}
}

func (s *SmtpSender) Send(to, subject, body string) error {
	addr := net.JoinHostPort(s.host, s.port)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.from, to, subject, body)

	return smtp.SendMail(addr, nil, s.from, []string{to}, []byte(msg))
}
