package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	sw "github.com/ZPI-2024-25/KubernetesAccessManager/api"
	"github.com/ZPI-2024-25/KubernetesAccessManager/auth"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	"github.com/ZPI-2024-25/KubernetesAccessManager/common"
	"github.com/ZPI-2024-25/KubernetesAccessManager/health"
	"github.com/gorilla/handlers"
)

func main() {
	requiredEnvVars := []string{
		"HEALTH_PORT",
		"APP_PORT",
		"VITE_KEYCLOAK_URL",
		"VITE_KEYCLOAK_CLIENTNAME",
		"VITE_KEYCLOAK_REALMNAME",
	}

	for _, envVar := range requiredEnvVars {
		common.GetEnvOrPanic(envVar)
	}

	healthPortStr := common.GetEnvOrPanic("HEALTH_PORT")
	healthPort, err := strconv.Atoi(healthPortStr)
	if err != nil {
		log.Fatalf("Invalid value for HEALTH_PORT: %s. Must be an integer. Exiting...", healthPortStr)
	}

	appPortStr := common.GetEnvOrPanic("APP_PORT")
	appPort, err := strconv.Atoi(appPortStr)
	if err != nil {
		log.Fatalf("Invalid value for APP_PORT: %s. Must be an integer. Exiting...", appPortStr)
	}

	healthServer := health.PrepareHealthEndpoints(healthPort)

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
		log.Printf("Health endpoints starting on port %d", healthPort)
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

	serverAddress := fmt.Sprintf(":%d", appPort)
	log.Printf("Starting server on %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, corsHandler(router)))
}
