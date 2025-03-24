package client

import (
	"context"

	"github.com/cloudnative-pg/cloudnative-pg/internal/cnpi/plugin/repository"
	"github.com/cloudnative-pg/cloudnative-pg/internal/configuration"
	"github.com/cloudnative-pg/machinery/pkg/log"
	"github.com/cloudnative-pg/machinery/pkg/stringset"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type getPlugins interface {
	getPlugins()
	client.Object
}

func NewClient(ctx context.Context, enabledPlugin *stringset.Data) (Client, error) {
	contextLogger := log.FromContext(ctx)
	plugins := repository.New()
	availablePluginNames, err := plugins.RegisterUnixSocketPluginsInPath(configuration.Current.PluginSocketDir)
	if err != nil {
		contextLogger.Error(err, "Error while loading local plugins")
		plugins.Close()
		return nil, err
	}

	availablePluginNamesSet := stringset.From(availablePluginNames)
	availableAndEnabled := stringset.From(availablePluginNamesSet.Intersect(enabledPlugin).ToList())
	return WithPlugins(
		ctx,
		plugins,
		availableAndEnabled.ToList()...,
	)
}
