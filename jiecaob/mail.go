package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"net"
	"strings"
	"time"
)

func SendMail(senderName string, sender string, to string, subject string, content string) error {
	// nslookup -q=mx qq.com
	mxs, err := net.LookupMX(to[strings.IndexRune(to, '@')+1:])
	if err != nil {
		return err
	}
	host := mxs[0].Host

	conn, err := net.DialTimeout("tcp", host+":25", 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	r := bufio.NewReader(conn)

	line, _, err := r.ReadLine()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(line), "220") {
		return errors.New(string(line))
	}
	conn.Write([]byte("HELO " + sender[strings.IndexRune(sender, '@')+1:] + "\r\n"))

	line, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(line), "250") {
		return errors.New(string(line))
	}
	conn.Write([]byte("MAIL FROM: <" + sender + ">\r\n"))

	line, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(line), "250") {
		return errors.New(string(line))
	}
	conn.Write([]byte("RCPT TO: <" + to + ">\r\n"))

	line, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(line), "250") {
		return errors.New(string(line))
	}
	conn.Write([]byte("DATA\r\n"))

	line, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(line), "354") {
		return errors.New(string(line))
	}
	var data = "From: =?utf8?B?" + toBase64(senderName) + "?= <" + sender + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: =?utf8?B?" + toBase64(subject) + "?=\r\n" +
		"Date: " + time.Now().Format(time.RFC3339) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"utf8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"X-Priority: 3\r\n" +
		"X-Mailer: golang Mail Sender\r\n" +
		"\r\n" +
		toBase64(content) +
		"\r\n.\r\n"
	conn.Write([]byte(data))

	line, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	if strings.HasPrefix(string(line), "550") {
		return errors.New("Mail is intercepted: " + string(line))
	}
	if !strings.HasPrefix(string(line), "250") {
		return errors.New(string(line))
	}
	conn.Write([]byte("QUIT\r\n"))

	line, _, err = r.ReadLine()
	if err != nil {
		return err
	}
	if !strings.HasPrefix(string(line), "221") {
		return errors.New(string(line))
	}

	return nil
}

func toBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
