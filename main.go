/*
Kaguya
===
Opensource websocket chat server for IM.

Copyright(c) 2020 Star Inc. All Rights Reserved.
The software licensed under Mozilla Public License Version 2.0
*/
package main

import (
	"fmt"
	"log"
	"net/http"

	kaguya "./libs"
	"github.com/gorilla/websocket"
)

func main() {
	kaguya.ReadConfig()
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer func() {
			log.Println("disconnect !!")
			c.Close()
		}()
		kaguya.HandleRequest(c)
	})

	fmt.Println("\n\tKaguya")
	fmt.Println("\t ===== ")
	fmt.Println("\n\tOpensource websocket chat server for IM.")
	fmt.Print("\n\tServer start at :8899\n\n")
	log.Fatal(http.ListenAndServe(":8899", nil))
	fmt.Print("\t\t(c) 2020 Star Inc. https://starinc.xyz\n\n")
}
