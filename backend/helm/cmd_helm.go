package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"time"
)

func getActionConfig(kubeConfig *rest.Config, namespace string) (*action.Configuration, error) {
	settings := cli.New()
	settings.KubeAPIServer = kubeConfig.Host
	settings.KubeToken = kubeConfig.BearerToken
	settings.KubeCaFile = kubeConfig.TLSClientConfig.CAFile
	settings.KubeInsecureSkipTLSVerify = kubeConfig.TLSClientConfig.Insecure
	settings.SetNamespace(namespace)

	configFlags := &genericclioptions.ConfigFlags{
		APIServer:   &kubeConfig.Host,
		CAFile:      &kubeConfig.TLSClientConfig.CAFile,
		BearerToken: &kubeConfig.BearerToken,
		Insecure:    &kubeConfig.TLSClientConfig.Insecure,
		Namespace:   &namespace,
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(configFlags, namespace, "secret", func(format string, v ...interface{}) {}); err != nil {
		return nil, err
	}

	return actionConfig, nil
}

func getRelease(actionConfig *action.Configuration, name string) (*release.Release, error) {
	get := action.NewGet(actionConfig)
	rel, err := get.Run(name)
	if err != nil {
		return nil, err
	}
	return rel, nil
}

func rollbackRelease(actionConfig *action.Configuration, name string, version int) error {
	rollback := action.NewRollback(actionConfig)
	rollback.Version = version
	rollback.Wait = true
	rollback.Timeout = 300 * time.Second
	if err := rollback.Run(name); err != nil {
		return err
	}

	return nil
}

func uninstallRelease(actionConfig *action.Configuration, name string) (*release.UninstallReleaseResponse, error) {
	uninstall := action.NewUninstall(actionConfig)
	response, err := uninstall.Run(name)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getReleaseHistory(actionConfig *action.Configuration, name string, max int) ([]*release.Release, error) {
	history := action.NewHistory(actionConfig)
	history.Max = max
	historyResponse, err := history.Run(name)
	if err != nil {
		return nil, err
	}

	return historyResponse, nil
}

func listReleases(actionConfig *action.Configuration, allNamespaces bool) ([]*release.Release, error) {
	list := action.NewList(actionConfig)
	if allNamespaces {
		list.AllNamespaces = true
	}
	list.StateMask = action.ListAll
	listResponse, err := list.Run()
	if err != nil {
		return nil, err
	}

	return listResponse, nil
}
