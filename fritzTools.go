package main

import (
	"fmt"
	"net/textproto"
	"strings"
)

func main() {
	conn, err := textproto.Dial("tcp", "router.home.lan:1012")
	defer conn.Close()
	if err != nil {
		panic(1)
	}
	fmt.Println("Verbindung zur FritzBox hergestellt")

	for {
		line, err := conn.Reader.ReadLine()
		if err != nil {
			panic(2)
		}

		callValues := strings.Split(line, ";")
		typeOfResp := callValues[1]

		switch typeOfResp {
		case "RING":
			fmt.Println("Anruf von " + callValues[3])
		case "CONNECT":
			fmt.Println("Angenommen von " + callValues[3])
		case "DISCONNECT":
			fmt.Println("Aufgelegt")
		}
	}
}
