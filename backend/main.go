package main

import (
	"fmt"
	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
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
	log.Fatal(http.ListenAndServe(":8080", router))
}
