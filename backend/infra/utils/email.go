package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type Email struct {
	To      string
	Subject string
	Body    string
}

const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

//SendEmailTLS send email on tls
func SendEmailTLS(config EmailConfig, e Email) (err error) {

	auth := smtp.PlainAuth("", config.Email, config.Password, config.Host)

	body := "From: " + "suporte@mvp.com.br" + "\nTo: " + e.To + "\r\nSubject: " + e.Subject + "\r\n" + mime + "\r\n" + e.Body
	err = smtp.SendMail(config.Host+":"+config.Port, auth, config.Email, []string{e.To}, []byte(body))
	return
}

//SendEmailSSL send email on ssl
func SendEmailSSL(config EmailConfig, e Email) (err error) {
	// Connect to the SMTP Server

	servername := fmt.Sprintf("%s:%s", config.Host, config.Port)

	auth := smtp.PlainAuth("", config.Email, config.Password, config.Host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.Host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return
	}
	defer conn.Close()

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = config.Email
	headers["To"] = e.To
	headers["Subject"] = e.Subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += mime + "\r\n" + e.Body

	c, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return
	}

	// To && From
	if err = c.Mail(config.Email); err != nil {
		return
	}

	if err = c.Rcpt(e.To); err != nil {
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return
	}

	err = w.Close()
	if err != nil {
		return
	}
	err = c.Quit()
	return
}
