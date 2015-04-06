package main

import (
	"fmt"
	"net/mail"
	"net/textproto"
	"net/smtp"
	"encoding/base64"
	"strings"
	"log"
	"code.google.com/p/gcfg"
)

type lineHandler struct {
    cn *textproto.Conn
}

func SendMail( subject string, mailtext string ) {
	var cfg MainConfig
    err := gcfg.ReadFileInto( &cfg, "/home/markus/gocode/src/github.com/madmax411/fritztools/fritzTools.gcfg")
    if err != nil {
        log.Fatal(err)
    }    

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
 
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	errmail := smtp.SendMail(
		smtpServer + ":" + cfg.SMTP.Port,
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
		//[]byte("This is the email body."),
	)
	if errmail != nil {
		log.Fatal(errmail)
	}
}

func (l *lineHandler) Watch() {
    lastAction := ""
	currAction := ""
    lastCallNo := ""

	for {
		line, err := l.cn.Reader.ReadLine()
		if err != nil {
			panic(2)
		}

		callValues := strings.Split(line, ";")
		currAction = callValues[1]

		switch currAction {
		case "RING":
			fmt.Println("Call from " + callValues[3])
			lastCallNo = callValues[3]
		case "CONNECT":
			fmt.Println("Connected with extention station #" + callValues[3])
		case "DISCONNECT":
			if lastAction == "RING" {
				fmt.Println("Send a info mail...")
				SendMail( "Fritz: Call", "Call from " + lastCallNo )
			}

			fmt.Println("Disconneted")
		}

		lastAction = currAction
	}
	
	return   
}

type MainConfig struct {
    Fritzbox Config_Fritzbox
    SMTP Config_SMTP
    Mail Config_Mail    
}

type Config_Fritzbox struct {
    Host string
    Port string
}

type Config_SMTP struct {
    Host string
    Port string
    User string
    Password string
}

type Config_Mail struct {
    From string
    To string
}




func main() {
	var cfg MainConfig
    err := gcfg.ReadFileInto( &cfg, "/home/markus/gocode/src/github.com/madmax411/fritztools/fritzTools.gcfg")
    if err != nil {
        log.Fatal(err)
    }    
   
    SendMail( "Fritz", "Test" )
   
    conn, err := textproto.Dial("tcp", cfg.Fritzbox.Host + ":" + cfg.Fritzbox.Port)
	defer conn.Close()
	if err != nil {
		panic(1)
	}
	
	fmt.Println("Connected to", cfg.Fritzbox.Host)
   
    lh := lineHandler{}
    lh.cn = conn

    lh.Watch()
}


