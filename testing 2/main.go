package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type sahamStruct struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan *sahamStruct)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("Go Websockets")
	router := mux.NewRouter()
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/data", dataHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandler)
	go publish()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func writer(saham *sahamStruct) {
	broadcast <- saham
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	var saham sahamStruct
	if err := json.NewDecoder(r.Body).Decode(&saham); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writer(&saham)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[ws] = true
}

func publish() {
	for {
		val := <-broadcast
		saham := fmt.Sprintf("%s %.3f", val.Name, val.Price)
		// send to every client that is currently connected
		fmt.Println(saham)
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(saham))
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
