package emailer

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// encode header of email message
func mailEncodeHeader(str string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(str)) + "?="
}

// encode header email field
func mailEncodeEmail(email string) string {
	return mailEncodeHeader(email) + " <" + email + ">"
}

// join encoded emails (receivers) for header field "To"
func mailJoinReceivers(receivers []string) string {
	arr := make([]string, len(receivers))
	for i, v := range receivers {
		arr[i] = mailEncodeEmail(v)
	}
	return strings.Join(arr, ",")
}

// send email with text/plain mime type
func SendEmail(subject, message string, receivers []string, userName, userPassword, host, identity string, port int16) (err error) {
	auth := smtp.PlainAuth(identity, userName, userPassword, host)
	msg := []byte("To: " + mailJoinReceivers(receivers) + "\r\n" +
		"Date:" + time.Now().Format("Mon 2 Jan 2006 15:04:05 -0700") + "\r\n" +
		"From: " + mailEncodeEmail(userName) + "\r\n" +
		"Subject: " + mailEncodeHeader(subject) + "\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" + message + "\r\n")
	return smtp.SendMail(fmt.Sprintf("%v:%v", host, port), auth, userName, receivers, msg)
}

// send email with text/html mime type
func SendEmailHTML(subject, message string, receivers []string, userName, userPassword, host, identity string, port int16) (err error) {
	auth := smtp.PlainAuth(identity, userName, userPassword, host)
	msg := []byte("To: " + mailJoinReceivers(receivers) + "\r\n" +
		"Date:" + time.Now().Format("Mon 2 Jan 2006 15:04:05 -0700") + "\r\n" +
		"From: " + mailEncodeEmail(userName) + "\r\n" +
		"Subject: " + mailEncodeHeader(subject) + "\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" + message + "\r\n")
	return smtp.SendMail(fmt.Sprintf("%v:%v", host, port), auth, userName, receivers, msg)
}
