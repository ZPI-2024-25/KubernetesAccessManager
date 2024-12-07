package main

import (
	"fmt"
	"log"
	"net/http"

	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"github.com/ZPI-2024-25/KubernetesAccessManager/auth"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/health"
	"github.com/gorilla/handlers"
)

func main() {
	common.InitEnv()
	auth.InitializeAuth()
	healthServer := health.PrepareHealthEndpoints(common.HealthPort)

	clusterSingleton, err := cluster.GetInstance()
	if err != nil {
		log.Fatalf("Error when loading config: %v\n", err)
	}

	_, err = auth.GetRoleMapInstance()
	if err != nil {
		log.Printf("Error when loading role map: %v\n", err)
	}

	go auth.WatchForRolemapChanges()
	go func() {
		log.Printf("Health endpoints starting on port %d", common.HealthPort)
		if err := healthServer.ListenAndServe(); err != nil {
			log.Fatal("Health endpoints have been shut down unexpectedly: ", err)
		}
	}()

	log.Printf("Marking application liveness as UP")
	health.ApplicationStatus.MarkAsUp()

	log.Printf("Server started")
	log.Printf("Authentication method: %s", clusterSingleton.GetAuthenticationMethod())

	router := sw.NewRouter()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)

	health.ServiceStatus.MarkAsUp()
	log.Printf("Marking application readiness as UP")

	serverAddress := fmt.Sprintf(":%d", common.AppPort)
	log.Printf("Starting server on %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, corsHandler(router)))
}
