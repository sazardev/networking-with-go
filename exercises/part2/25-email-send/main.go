// 25-email-send/main.go
// Builds a multipart/mixed MIME message (plain-text body plus a small
// attachment) and sends it with net/smtp. Replace the host/credentials
// below with a real (or local test, e.g. MailHog / Mailpit) SMTP server
// to actually deliver it -- as configured this will print a connection
// error, since smtp.example.com isn't a real server.
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
)

func buildMessage(
	from, to, subject, bodyText, filename string, attachment []byte,
) []byte {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fmt.Fprintf(&buf, "From: %s\r\n", from)
	fmt.Fprintf(&buf, "To: %s\r\n", to)
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf,
		"Content-Type: multipart/mixed; boundary=%q\r\n\r\n",
		writer.Boundary())

	textPart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/plain; charset=utf-8"},
	})
	textPart.Write([]byte(bodyText))

	attachPart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"application/octet-stream"},
		"Content-Transfer-Encoding": {"base64"},
		"Content-Disposition": {
			fmt.Sprintf(`attachment; filename=%q`, filename),
		},
	})
	encoded := make(
		[]byte, base64.StdEncoding.EncodedLen(len(attachment)))
	base64.StdEncoding.Encode(encoded, attachment)
	attachPart.Write(encoded)

	writer.Close()
	return buf.Bytes()
}

func main() {
	from := "sender@example.com"
	to := "recipient@example.com"

	msg := buildMessage(
		from, to, "Report attached",
		"Here's the file you asked for.\r\n",
		"notes.txt", []byte("just a small test attachment"))

	fmt.Println(string(msg))

	auth := smtp.PlainAuth("", from, "app-password", "smtp.example.com")
	err := smtp.SendMail(
		"smtp.example.com:587", auth, from, []string{to}, msg)
	if err != nil {
		fmt.Println("send error:", err)
		return
	}
	fmt.Println("sent")
}
