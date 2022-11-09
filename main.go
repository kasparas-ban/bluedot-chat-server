package main

import (
	env "bluedot-chat-server/config"
	db "bluedot-chat-server/database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
)

func loadEnv() {
	env.SERVER_PORT = os.Getenv("SERVER_PORT")
	env.ALLOWED_ORIGIN = os.Getenv("ALLOWED_ORIGIN")
	env.CONNECTIONS_DB_PASS = os.Getenv("CONNECTIONS_DB_PASS")
}

func main() {
	loadEnv()

	db.ConnectionsDB.Connect(fmt.Sprintf("redis://default:%s@connections_db:6379/0", env.CONNECTIONS_DB_PASS))

	http.HandleFunc("/", ws)
	if err := http.ListenAndServe(":"+env.SERVER_PORT, nil); err != nil {
		log.Fatal(err)
	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	// Upgrade connection
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == env.ALLOWED_ORIGIN
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to establish connection", http.StatusInternalServerError)
		return
	}

	// Add user to connectionsDB
	err = db.AddConn(username)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// Read messages from socket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, "failed to read the message from the client. Ending the connection", http.StatusInternalServerError)
			conn.Close()
			return
		}
		log.Printf("msg: %s", string(msg))
	}
}
