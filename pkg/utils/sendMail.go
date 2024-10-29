package utils

import (
	"log"
	"net/smtp"
)

func SendMail(to []string, subject, body string) {
    from := ""
    appPassword := ""

    host := "smtp.gmail.com"
    port := "587"
    address := host+":"+port

    message := []byte(subject+body)
    auth := smtp.PlainAuth("", from, appPassword, host)

    err := smtp.SendMail(address, auth, from, to, message)

    if err != nil {
        log.Printf("Failed to send the email: %v", err)
    }
}

