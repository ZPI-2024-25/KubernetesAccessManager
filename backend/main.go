package main

import (
	"flag"
	"fmt"
	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"log"
	"net/http"
)

func main() {
	flag.Parse()

	_, err := common.GetInstance()
	if err != nil {
		fmt.Printf("Error when loading config: %v\n", err)
		return
	}

	log.Printf("Server started")
	router := sw.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
