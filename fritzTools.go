package main

import (
	"fmt"
	"net/textproto"
	"strings"
)

func main() {
	host := "router.home.lan"
	port := "1012"

	conn, err := textproto.Dial("tcp", host+":"+port)
	defer conn.Close()
	if err != nil {
		panic(1)
	}
	fmt.Println("Connected to", host)

	lastAction := ""
	currAction := ""

	for {
		line, err := conn.Reader.ReadLine()
		if err != nil {
			panic(2)
		}

		callValues := strings.Split(line, ";")
		currAction = callValues[1]

		switch currAction {
		case "RING":
			fmt.Println("Call from " + callValues[3])
		case "CONNECT":
			fmt.Println("Connected with extention station #" + callValues[3])
		case "DISCONNECT":
			if lastAction == "RING" {
				fmt.Println("TODO: Send a info mail...")
			}

			fmt.Println("Disconneted")
		}

		lastAction = currAction
	}
}
