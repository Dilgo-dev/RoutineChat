package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dilgo-dev/RoutineChat/internal/handlers"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/template/index.html")
	})

	http.Handle("/ws", websocket.Handler(handlers.NewServer().HandleWS))

	fmt.Println("Server is running on port http://localhost:" + port)

	http.ListenAndServe(":"+port, nil)
}
