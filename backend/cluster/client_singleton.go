package cluster

import (
	"flag"
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sync"

	"k8s.io/client-go/rest"
)

type DynamicClientSingleton struct {
	config        *rest.Config
	dynamicClient *dynamic.DynamicClient
}

var (
	instance      *DynamicClientSingleton
	once          sync.Once
	kubeconfig    *string
	inClusterAuth *bool
)

func init() {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	inClusterAuth = flag.Bool("in-cluster", false, "Use in-cluster authentication")
}

func GetInstance() (*DynamicClientSingleton, error) {
	var err error
	once.Do(func() {
		var config *rest.Config

		if *inClusterAuth {
			config, err = rest.InClusterConfig()
			if err != nil {
				err = fmt.Errorf("failed to get in-cluster config: %w", err)
				return
			}
		} else {
			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
			if err != nil {
				err = fmt.Errorf("nie udało się zbudować konfiguracji z pliku kubeconfig: %w", err)
				return
			}
		}

		dynamicClient, innerErr := dynamic.NewForConfig(config)
		if innerErr != nil {
			err = fmt.Errorf("nie udało się utworzyć dynamicznego klienta: %w", innerErr)
			return
		}

		instance = &DynamicClientSingleton{
			config:        config,
			dynamicClient: dynamicClient,
		}
	})

	if instance == nil {
		return nil, err
	}

	return instance, nil
}

func (c *DynamicClientSingleton) GetClientSet() *dynamic.DynamicClient {
	return c.dynamicClient
}

func (c *DynamicClientSingleton) GetConfig() *rest.Config {
	return c.config
}

func (c *DynamicClientSingleton) GetAuthenticationMethod() string {
	if *inClusterAuth {
		return "in-cluster"
	} else {
		return "out-of-cluster"
	}
}
