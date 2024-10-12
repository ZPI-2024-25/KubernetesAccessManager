package common

import (
	"fmt"
	"sync"

	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ClientSetSingleton struct {
	config    *rest.Config
	clientset *kubernetes.Clientset
}

var (
	instance *ClientSetSingleton
	once     sync.Once
)

func GetInstance() (*ClientSetSingleton, error) {
	var err error
	once.Do(func() {
		config, innerErr := rest.InClusterConfig()
		if err != nil {
			err = fmt.Errorf("failed to get in-cluster config: %w", innerErr)
			return
		}

		clientset, innerErr := kubernetes.NewForConfig(config)

		if err != nil {
			err = fmt.Errorf("failed to create clientset: %w", innerErr)
			return
		}

		instance = &ClientSetSingleton{
			config:    config,
			clientset: clientset,
		}
	})

	if instance == nil {
		return nil, err
	}

	return instance, nil
}

func (s *ClientSetSingleton) QueryPods(namespace string) ([]string, error) {
	pods, err := s.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("pods not found")
		} else {
			return nil, fmt.Errorf("error: %s", err.Error())
		}
	}

	var podNames []string

	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

func (s *ClientSetSingleton) QueryDeployments(namespace string) ([]string, error) {
	deployments, err := s.clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("deployments not found")
		} else {
			return nil, fmt.Errorf("error: %s", err.Error())
		}
	}

	var deploymentNames []string

	for _, deployment := range deployments.Items {
		deploymentNames = append(deploymentNames, deployment.Name)
	}

	return deploymentNames, nil
}

func (s *ClientSetSingleton) QueryServices(namespace string) ([]string, error) {
	services, err := s.clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("services not found")
		} else {
			return nil, fmt.Errorf("error: %s", err.Error())
		}
	}

	var serviceNames []string

	for _, service := range services.Items {
		serviceNames = append(serviceNames, service.Name)
	}

	return serviceNames, nil
}
