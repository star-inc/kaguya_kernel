/*
Kaguya
===
The opensource instant messaging framework.

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

func declare(port string) {
	fmt.Println("\n\tKaguya")
	fmt.Println("\t ===== ")
	fmt.Println("\n\tOpensource websocket chat server for IM.")
	fmt.Printf("\n\tServer start at %s\n\n", port)
	fmt.Print("\t\t(c) 2020 Star Inc. https://starinc.xyz\n\n")
}

func main() {
	kaguya.ReadConfig()
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() {
			c.Close()
		}()
		kaguya.NewHandleInterface(c).Start()
	})
	servePort := fmt.Sprintf(":%d", kaguya.Config.Port)
	declare(servePort)
	log.Fatal(http.ListenAndServe(servePort, nil))
}
