package main

import (
	"fmt"
	"net/smtp"
)

// Example: Send an email using Go's net/smtp package.
func main() {
	from := "your@email.com"
	password := "yourpassword"
	to := []string{"recipient@email.com"}
	host := "smtp.example.com"
	port := "587"
	msg := []byte("Subject: Hello from Go!\r\n\r\nThis is a test email sent from Go.")

	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}
	fmt.Println("Email sent!")
}
