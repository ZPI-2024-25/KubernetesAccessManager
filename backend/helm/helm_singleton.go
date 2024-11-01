package helm

import (
	"fmt"
	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	helmclient "github.com/mittwald/go-helm-client"
	"k8s.io/client-go/rest"
	"sync"
)

type HelmClientSingleton struct {
	helmClient helmclient.Client
}

var (
	instance *HelmClientSingleton
	once     sync.Once
)

func GetInstance() (*HelmClientSingleton, error) {
	var err error
	once.Do(func() {
		var config *rest.Config
		config, err = cluster.GetConfig()
		if err != nil {
			return
		}

		opt := &helmclient.RestConfClientOptions{
			Options: &helmclient.Options{
				Namespace: "",
				Debug:     true,
				Linting:   true,
			},
			RestConfig: config,
		}

		helmClient, innerErr := helmclient.NewClientFromRestConf(opt)

		if innerErr != nil {
			err = fmt.Errorf("failed to create helm client: %w", innerErr)
			return
		}

		instance = &HelmClientSingleton{
			helmClient: helmClient,
		}
	})

	if instance == nil {
		return nil, err
	}

	return instance, nil
}

func GetHelmClient() (helmclient.Client, error) {
	singleton, err := GetInstance()
	if err != nil {
		return nil, err
	}

	return singleton.helmClient, nil
}
