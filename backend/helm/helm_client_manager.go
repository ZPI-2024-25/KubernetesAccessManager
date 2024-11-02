package helm

import (
	"fmt"
	"sync"

	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	helmclient "github.com/mittwald/go-helm-client"
	"k8s.io/client-go/rest"
)

type HelmClientManager struct {
	clients map[string]helmclient.Client
	config  *rest.Config
	mu      sync.Mutex
}

var (
	instance *HelmClientManager
	once     sync.Once
)

func GetInstance() (*HelmClientManager, error) {
	var err error
	once.Do(func() {
		var config *rest.Config
		config, err = cluster.GetConfig()
		if err != nil {
			return
		}

		instance = &HelmClientManager{
			clients: make(map[string]helmclient.Client),
			config:  config,
		}
	})

	if instance == nil {
		return nil, err
	}

	return instance, nil
}

func GetHelmClient(namespace string) (helmclient.Client, error) {
	manager, err := GetInstance()
	if err != nil {
		return nil, err
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	if client, exists := manager.clients[namespace]; exists {
		return client, nil
	}

	opt := &helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace: namespace,
			Debug:     true,
			Linting:   true,
		},
		RestConfig: manager.config,
	}

	helmClient, err := helmclient.NewClientFromRestConf(opt)
	if err != nil {
		return nil, fmt.Errorf("nie udało się utworzyć klienta Helm dla namespace '%s': %w", namespace, err)
	}

	manager.clients[namespace] = helmClient

	return helmClient, nil
}
