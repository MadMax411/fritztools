package main

import (
	"code.google.com/p/gcfg"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
)

type lineHandler struct {
	cn  *textproto.Conn
	cfg MainConfig
}

type lastCall struct {
	PhoneNo  string
	DateTime string
	ToNo     string
}

type MainConfig struct {
	Fritzbox Config_Fritzbox
	SMTP     Config_SMTP
	Mail     Config_Mail
}

type Config_Fritzbox struct {
	Host string
	Port string
}

type Config_SMTP struct {
	SendMail bool
	Host     string
	Port     string
	User     string
	Password string
}

type Config_Mail struct {
	From string
	To   string
}

func SendMail(subject string, mailtext string, cfg MainConfig) {

	if cfg.SMTP.SendMail {
		smtpServer := cfg.SMTP.Host
		auth := smtp.PlainAuth(
			"",
			cfg.SMTP.User,
			cfg.SMTP.Password,
			smtpServer,
		)

		from := mail.Address{"", cfg.Mail.From}
		to := mail.Address{"", cfg.Mail.To}
		title := subject

		body := mailtext

		header := make(map[string]string)
		header["From"] = from.String()
		header["To"] = to.String()
		header["Subject"] = title
		header["MIME-Version"] = "1.0"
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
		header["Content-Transfer-Encoding"] = "base64"

		message := ""
		for k, v := range header {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

		errmail := smtp.SendMail(
			smtpServer+":"+cfg.SMTP.Port,
			auth,
			from.Address,
			[]string{to.Address},
			[]byte(message),
		)
		if errmail != nil {
			log.Fatal(errmail)
		}
	}
}

func (l *lineHandler) Watch() {
	lastAction := ""
	currAction := ""
	extStation := ""

	call := lastCall{}

	for {
		line, err := l.cn.Reader.ReadLine()
		if err != nil {
			panic(err)
		}

		callValues := strings.Split(line, ";")
		currAction = callValues[1]

		switch currAction {
		case "RING":
			//fmt.Println("Call from " + callValues[3] + " to " + callValues[4])
			call.PhoneNo = callValues[3]
			call.DateTime = callValues[0]
			call.ToNo = callValues[4]

		case "CONNECT":
			//fmt.Println("Connected with extention station #" + callValues[3])
			extStation = callValues[3]

		case "DISCONNECT":
			msg := ""

			if lastAction == "RING" {
				//fmt.Print("Send a info mail...")
				msg = "Unanswered call from " + call.PhoneNo + " at " + call.DateTime
				msg += "\nTo number " + call.ToNo
			}

			if lastAction == "CONNECT" && extStation == "40" {
				//fmt.Print("Send a info mail...")
				msg = "Call from " + call.PhoneNo + " at " + call.DateTime
				msg += "\nTo number " + call.ToNo
				msg += "\nCall is answered by the answering machine"
			}

			SendMail("Fritz: PhoneCall", msg, l.cfg)
			//fmt.Println("Disconneted")
		}

		lastAction = currAction
	}

	return
}

func main() {
	var cfg MainConfig
	err := gcfg.ReadFileInto(&cfg, "fritzTools.gcfg")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	conn, err := textproto.Dial("tcp", cfg.Fritzbox.Host+":"+cfg.Fritzbox.Port)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	lh := lineHandler{}
	lh.cn = conn
	lh.cfg = cfg

	lh.Watch()
}
