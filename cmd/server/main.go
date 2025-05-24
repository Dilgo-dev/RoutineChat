package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dilgo-dev/RoutineChat/internal/config"
	"github.com/Dilgo-dev/RoutineChat/internal/handlers"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/template/index.html")
	})

	http.Handle("/ws", websocket.Handler(handlers.NewServer().HandleWS))

	fmt.Println("Server is running on port http://localhost:" + cfg.Port)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
