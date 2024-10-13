package main

import (
	sw "github.com/ZPI-2024-25/KubernetesUserManager/api"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")
	router := sw.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
	TODO
    Deploy folder and Deployment Instructions.md should be removed  when logic under api endpoints
    is implemented.
*/
