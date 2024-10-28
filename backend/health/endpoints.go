package health

import (
	"fmt"
	"github.com/Icikowski/kubeprobes"
	"net/http"
)

func PrepareHealthEndpoints(port int) *http.Server {

	health := kubeprobes.New(
		kubeprobes.WithLivenessProbes(ApplicationStatus.GetProbeFunction()),
		kubeprobes.WithReadinessProbes(ServiceStatus.GetProbeFunction()),
	)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: health,
	}
}
