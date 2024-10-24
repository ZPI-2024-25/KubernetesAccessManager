package main

import (
	"fmt"
	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/health"
	"log"
	"net/http"
)

func main() {
	healthServer := health.PrepareHealthEndpoints(
		8082,
	)
	singleton, err := cluster.GetInstance()
	if err != nil {
		fmt.Printf("Error when loading config: %v\n", err)
		return
	}
	go func() {
		log.Printf("health endpoints starting")
		if err := healthServer.ListenAndServe(); err != nil {
			log.Fatal("health endpoints have been shut down unexpectedly: ", err)
		}
	}()
	log.Printf("marking application liveness as UP")
	health.ApplicationStatus.MarkAsUp()
	log.Printf("Server started")
	log.Printf("Authentication method: %s", singleton.GetAuthenticationMethod())
	router := sw.NewRouter()
	health.ServiceStatus.MarkAsUp()
	log.Printf("marking application readiness as UP")
	log.Fatal(http.ListenAndServe(":8080", router))
}
