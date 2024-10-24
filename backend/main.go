package main

import (
	"fmt"
	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {
	singleton, err := cluster.GetInstance()
	if err != nil {
		fmt.Printf("Error when loading config: %v\n", err)
		return
	}

	log.Printf("Server started")
	log.Printf("Authentication method: %s", singleton.GetAuthenticationMethod())

	router := sw.NewRouter()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)

	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}
