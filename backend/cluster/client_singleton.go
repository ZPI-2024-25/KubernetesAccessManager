package cluster

import (
	"flag"
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
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
	kubeconfig    string
	inClusterAuth bool
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	flag.BoolVar(&inClusterAuth, "in-cluster", false, "Use in-cluster authentication")
}

func GetInstance() (*DynamicClientSingleton, error) {
	var err error
	once.Do(func() {
		var config *rest.Config

		flag.Parse()

		if inClusterAuth {
			config, err = rest.InClusterConfig()
			if err != nil {
				log.Printf("Error when loading in-cluster config: %v\n", err)
				return
			}
		} else {
			kubeconfigPath := kubeconfig

			if kubeconfigPath == "" {
				kubeconfigEnv := os.Getenv("KUBECONFIG")
				if kubeconfigEnv != "" {
					kubeconfigPath = kubeconfigEnv
				} else {
					if home := homedir.HomeDir(); home != "" {
						kubeconfigPath = filepath.Join(home, ".kube", "config")
					} else {
						err = fmt.Errorf("could not determine home directory")
						return
					}
				}
			}

			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
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
	if inClusterAuth {
		return "in-cluster"
	} else {
		return "out-of-cluster"
	}
}
