package helm

import (
	"fmt"
	"sync"

	"github.com/ZPI-2024-25/KubernetesAccessManager/cluster"
	helmclient "github.com/mittwald/go-helm-client"
	"k8s.io/client-go/rest"
)

type HelmClientSingleton struct {
	helmClient helmclient.Client
	namespace  string
	config     *rest.Config
	mu         sync.Mutex
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
			namespace:  "",
			config:     config,
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

	singleton.mu.Lock()
	defer singleton.mu.Unlock()

	return singleton.helmClient, nil
}

func GetCurrentNamespace() string {
	instance.mu.Lock()
	defer instance.mu.Unlock()
	return instance.namespace
}

func RegenerateWithNewNamespace(namespace string) error {
	singleton, err := GetInstance()
	if err != nil {
		return err
	}

	singleton.mu.Lock()
	defer singleton.mu.Unlock()

	if namespace != singleton.namespace {
		opt := &helmclient.RestConfClientOptions{
			Options: &helmclient.Options{
				Namespace: namespace,
				Debug:     true,
				Linting:   true,
			},
			RestConfig: singleton.config,
		}

		helmClient, innerErr := helmclient.NewClientFromRestConf(opt)
		if innerErr != nil {
			return fmt.Errorf("failed to create helm client with new namespace: %w", innerErr)
		}

		singleton.helmClient = helmClient
		singleton.namespace = namespace
	}

	return nil
}
