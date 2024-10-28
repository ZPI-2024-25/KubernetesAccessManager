package main

import (
	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")
	router := sw.NewRouter()
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)

	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}
