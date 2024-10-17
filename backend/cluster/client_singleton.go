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
				err = fmt.Errorf("failed to create config from kubeconfig: %w", err)
				return
			}
		}

		dynamicClient, innerErr := dynamic.NewForConfig(config)
		if innerErr != nil {
			err = fmt.Errorf("failed to create dynamic client: %w", innerErr)
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

func GetClientSet() (*dynamic.DynamicClient, error) {
	singleton, err := GetInstance()
	if err != nil {
		return nil, err
	}

	return singleton.dynamicClient, nil
}

func GetConfig() (*rest.Config, error) {
	singleton, err := GetInstance()
	if err != nil {
		return nil, err
	}

	return singleton.config, nil
}

func (c *DynamicClientSingleton) GetAuthenticationMethod() string {
	if *inClusterAuth {
		return "in-cluster"
	} else {
		return "out-of-cluster"
	}
}
