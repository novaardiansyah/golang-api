package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

type MailData struct {
	Title           string
	Header          string
	Content         string
	AuthorName      string
	AuthorFirstName string
	Year            int
}

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
}

func SendEmail(conf SMTPConfig, to string, subject string, templatePath string, data interface{}) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", conf.Username, conf.Password, conf.Host)

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", conf.FromName, conf.FromAddress)
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	return smtp.SendMail(addr, auth, conf.FromAddress, []string{to}, []byte(message))
}

func GetDefaultMailData(authorName, header, content string) MailData {
	return MailData{
		Title:           header,
		Header:          header,
		Content:         content,
		AuthorName:      authorName,
		AuthorFirstName: "Nova",
		Year:            time.Now().Year(),
	}
}
