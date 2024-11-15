package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"time"
)

type ActionConfigInterface interface {
	getRelease(name string) (*release.Release, error)
	listReleases(allNamespaces bool) ([]*release.Release, error)
	uninstallRelease(name string) (*release.UninstallReleaseResponse, error)
	getReleaseHistory(name string, max int) ([]*release.Release, error)
	rollbackRelease(name string, version int) error
}

type ActionConfig struct {
	config *action.Configuration
}

func getActionConfig(kubeConfig *rest.Config, namespace string) (*ActionConfig, error) {
	configFlags := &genericclioptions.ConfigFlags{
		Namespace: &namespace,
		WrapConfigFn: func(_ *rest.Config) *rest.Config {
			return kubeConfig
		},
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(configFlags, namespace, "secret", func(format string, v ...interface{}) {}); err != nil {
		return nil, err
	}

	return &ActionConfig{config: actionConfig}, nil
}

func (c *ActionConfig) getRelease(name string) (*release.Release, error) {
	get := action.NewGet(c.config)
	rel, err := get.Run(name)
	if err != nil {
		return nil, err
	}
	return rel, nil
}

func (c *ActionConfig) rollbackRelease(name string, version int) error {
	rollback := action.NewRollback(c.config)
	rollback.Version = version
	rollback.Wait = true
	rollback.Timeout = 300 * time.Second
	if err := rollback.Run(name); err != nil {
		return err
	}

	return nil
}

func (c *ActionConfig) uninstallRelease(name string) (*release.UninstallReleaseResponse, error) {
	uninstall := action.NewUninstall(c.config)
	response, err := uninstall.Run(name)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ActionConfig) getReleaseHistory(name string, max int) ([]*release.Release, error) {
	history := action.NewHistory(c.config)
	history.Max = max
	historyResponse, err := history.Run(name)
	if err != nil {
		return nil, err
	}

	return historyResponse, nil
}

func (c *ActionConfig) listReleases(allNamespaces bool) ([]*release.Release, error) {
	list := action.NewList(c.config)
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
