// 25-email-parse/main.go
// Parses a raw RFC 5322 message with net/mail (headers, address list),
// then decodes its multipart/mixed body with mime and mime/multipart --
// the two packages net/mail deliberately leaves to you.
package main

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

const raw = "From: \"Ada Lovelace\" <ada@example.com>\r\n" +
	"To: alan@example.com\r\n" +
	"Subject: Re: the analytical engine\r\n" +
	"MIME-Version: 1.0\r\n" +
	"Content-Type: multipart/mixed; boundary=\"sep123\"\r\n" +
	"\r\n" +
	"--sep123\r\n" +
	"Content-Type: text/plain; charset=utf-8\r\n" +
	"\r\n" +
	"Looking forward to your notes.\r\n" +
	"--sep123\r\n" +
	"Content-Type: application/octet-stream\r\n" +
	"Content-Disposition: attachment; filename=\"notes.txt\"\r\n" +
	"\r\n" +
	"(pretend this is base64-encoded file data)\r\n" +
	"--sep123--\r\n"

func main() {
	msg, err := mail.ReadMessage(strings.NewReader(raw))
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	from, err := mail.ParseAddress(msg.Header.Get("From"))
	if err == nil {
		fmt.Printf("From: %s <%s>\n", from.Name, from.Address)
	}
	fmt.Println("Subject:", msg.Header.Get("Subject"))

	mediaType, params, err := mime.ParseMediaType(
		msg.Header.Get("Content-Type"))
	if err != nil {
		fmt.Println("content-type error:", err)
		return
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		body, _ := io.ReadAll(msg.Body)
		fmt.Println("Body:", string(body))
		return
	}

	reader := multipart.NewReader(msg.Body, params["boundary"])
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("part error:", err)
			return
		}
		data, _ := io.ReadAll(part)
		fmt.Printf("part %q: %d bytes -> %q\n",
			part.Header.Get("Content-Type"), len(data), data)
	}
}
