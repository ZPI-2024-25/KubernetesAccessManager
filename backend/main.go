package main

import (
	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")
	router := sw.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
