package cluster

import (
	"context"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes/secrets"
	"go.uber.org/zap"
)

type SecretLoader struct {
	secrets secrets.Client
	logging bool
}

func NewSecretLoader(secrets secrets.Client, logging bool) *SecretLoader {
	return &SecretLoader{
		secrets: secrets,
		logging: logging,
	}
}

func (l *SecretLoader) LoadConfigs(ctx context.Context, clusters []Config) []Config {
	loaded := make([]Config, 0, len(clusters))
	for _, cl := range clusters {
		cl, err := l.resolveCluster(ctx, cl)
		if err != nil {
			zap.L().Error("failed to load cluster secret", zap.Error(err), zap.String("cluser", cl.Name), zap.String("secretRef", cl.SecretRef))
			continue
		}
		cl.Logging = l.logging

		for i, p := range cl.Plugins {
			p, err := l.resolvePlugin(ctx, p)
			if err != nil {
				zap.L().Error(
					"failed to load plugin secret",
					zap.Error(err),
					zap.String("cluster", cl.Name),
					zap.String("plugin", p.Name),
					zap.String("secretRef", p.SecretRef),
				)
				continue
			}
			p.Logging = l.logging

			cl.Plugins[i] = p
		}

		loaded = append(loaded, cl)
	}

	return loaded
}

func (l *SecretLoader) LoadConfig(ctx context.Context, cluster Config) Config {
	cl, err := l.resolveCluster(ctx, cluster)
	if err != nil {
		zap.L().Error("failed to load cluster secret", zap.Error(err), zap.String("cluser", cl.Name), zap.String("secretRef", cl.SecretRef))
		return cluster
	}
	cl.Logging = l.logging

	for i, p := range cl.Plugins {
		p, err := l.resolvePlugin(ctx, p)
		if err != nil {
			zap.L().Error(
				"failed to load plugin secret",
				zap.Error(err),
				zap.String("cluster", cl.Name),
				zap.String("plugin", p.Name),
				zap.String("secretRef", p.SecretRef),
			)
			continue
		}
		p.Logging = l.logging
		cl.Plugins[i] = p
	}

	return cl
}

func (l *SecretLoader) resolveCluster(ctx context.Context, cl Config) (Config, error) {
	if cl.SecretRef != "" && l.secrets != nil {
		values, err := l.secrets.Get(ctx, cl.SecretRef)
		if err != nil {
			return cl, err
		}

		cl = cl.FromValues(values)
	}

	return cl, nil
}

func (l *SecretLoader) resolvePlugin(ctx context.Context, plugin Plugin) (Plugin, error) {
	if plugin.SecretRef != "" && l.secrets != nil {
		values, err := l.secrets.Get(ctx, plugin.SecretRef)
		if err != nil {
			return plugin, err
		}

		plugin = plugin.FromValues(values)
	}

	return plugin, nil
}
