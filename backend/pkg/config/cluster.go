package config

import (
	"context"
	"net/http/httputil"
	"net/url"

	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/api/proxy"
	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"github.com/kyverno/policy-reporter-ui/pkg/server"
	"go.uber.org/zap"
)

type ClusterOptions struct {
	Logging       bool
	OverwriteHost bool
}

type ClusterManager struct {
	client  secrets.Client
	server  server.Server
	options ClusterOptions
}

func (c *ClusterManager) Register(ctx context.Context, cluster Cluster) error {
	cluster, err := c.loadClusterSecret(ctx, cluster)
	if err != nil {
		zap.L().Error("failed to load cluster secret", zap.Error(err), zap.String("cluser", cluster.Name), zap.String("secretRef", cluster.SecretRef))
		return err
	}

	proxy, err := c.proxy(cluster)
	if err != nil {
		zap.L().Error("failed to resolve proxies", zap.Error(err), zap.String("cluser", cluster.Name), zap.String("host", cluster.Host))
		return err
	}

	client, err := c.coreClient(cluster)

	plugins := make(map[string]*plugin.Client, len(cluster.Plugins))
	for _, p := range cluster.Plugins {
		p, err := c.loadPluginSecret(ctx, p)
		if err != nil {
			zap.L().Error(
				"failed to load plugin secret",
				zap.Error(err),
				zap.String("cluster", cluster.Name),
				zap.String("plugin", p.Name),
				zap.String("secretRef", p.SecretRef),
			)
			continue
		}

		pClient, err := c.pluginClient(p)
		if err != nil {
			zap.L().Error("failed to create plugin client", zap.Error(err), zap.String("cluser", cluster.Name), zap.String("plugin", p.Name))
			continue
		}

		plugins[p.Name] = pClient
	}

	c.server.RegisterCluster(cluster.Name, client, plugins, proxy)

	return nil
}

func (c *ClusterManager) Watch(ctx context.Context) error {
	events, err := c.client.Watch(ctx)
	if err != nil {
		return err
	}

	for event := range events {
		switch event.Type {
		case secrets.Added:
			c.Register(ctx, (Cluster{}).FromValues(event.Values))
		}
	}

	return nil
}

func (c *ClusterManager) loadClusterSecret(ctx context.Context, cluster Cluster) (Cluster, error) {
	if cluster.SecretRef == "" {
		return cluster, nil
	}

	values, err := c.client.Get(ctx, cluster.SecretRef)
	if err != nil {
		return cluster, err
	}

	return cluster.FromValues(values), nil
}

func (c *ClusterManager) loadPluginSecret(ctx context.Context, plugin Plugin) (Plugin, error) {
	if plugin.SecretRef == "" {
		return plugin, nil
	}

	values, err := c.client.Get(ctx, plugin.SecretRef)
	if err != nil {
		return plugin, err
	}

	return plugin.FromValues(values), nil
}

func (c *ClusterManager) coreClient(cluster Cluster) (*core.Client, error) {
	options := []api.ClientOption{
		api.WithBaseURL(cluster.Host),
	}

	if cluster.Certificate != "" {
		options = append(options, api.WithCertificate(cluster.Certificate))
	} else if cluster.SkipTLS {
		options = append(options, api.WithSkipTLS())
	}

	if cluster.BasicAuth.Username != "" {
		options = append(options, api.WithBaseAuth(api.BasicAuth{
			Username: cluster.BasicAuth.Username,
			Password: cluster.BasicAuth.Password,
		}))
	}

	if c.options.Logging {
		options = append(options, api.WithLogging())
	}

	return core.New(options)
}

func (c *ClusterManager) pluginClient(p Plugin) (*plugin.Client, error) {
	options := []api.ClientOption{
		api.WithBaseURL(p.Host),
	}

	if p.Certificate != "" {
		options = append(options, api.WithCertificate(p.Certificate))
	} else if p.SkipTLS {
		options = append(options, api.WithSkipTLS())
	}

	if p.BasicAuth.Username != "" {
		options = append(options, api.WithBaseAuth(api.BasicAuth{
			Username: p.BasicAuth.Username,
			Password: p.BasicAuth.Password,
		}))
	}

	if c.options.Logging {
		options = append(options, api.WithLogging())
	}

	return plugin.New(options)
}

func (c *ClusterManager) proxy(cluster Cluster) (*httputil.ReverseProxy, error) {
	if cluster.Host == "" {
		return nil, ErrMissingAPI
	}

	target, err := url.Parse(cluster.Host)
	if err != nil {
		return nil, err
	}

	options := make([]proxy.DirectorOption, 0)
	proxyOptions := make([]proxy.ProxyOption, 0)
	basicAuth := cluster.BasicAuth

	if c.options.Logging {
		options = append(options, proxy.WithLogging())
	}

	if c.options.OverwriteHost {
		options = append(options, proxy.WithHostOverwrite())
	}

	if basicAuth.Username != "" && basicAuth.Password != "" {
		options = append(options, proxy.WithAuth(basicAuth.Username, basicAuth.Password))
	}

	if cluster.SkipTLS {
		proxyOptions = append(proxyOptions, proxy.WithSkipTLS())
	}

	if cluster.Certificate != "" {
		proxyOptions = append(proxyOptions, proxy.WithCertificate(cluster.Certificate))
	}

	return proxy.New(target, options, proxyOptions), nil
}

func (c *ClusterManager) LoadBasicAuth(ctx context.Context, secretRef string) (*BasicAuth, error) {
	values, err := c.client.Get(ctx, secretRef)
	if err != nil {
		return nil, err
	}

	return &BasicAuth{
		Username:  values.Username,
		Password:  values.Password,
		SecretRef: secretRef,
	}, nil
}
