package health

import (
	"github.com/Icikowski/kubeprobes"
)

var (
	ApplicationStatus = kubeprobes.NewStatefulProbe()
	ServiceStatus     = kubeprobes.NewStatefulProbe()
)
