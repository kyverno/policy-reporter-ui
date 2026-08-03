package cluster

import (
	"github.com/kyverno/policy-reporter-ui/pkg/api"
	"github.com/kyverno/policy-reporter-ui/pkg/api/core"
	"github.com/kyverno/policy-reporter-ui/pkg/api/plugin"
	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/ui/v1alpha1"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	"go.uber.org/zap"
)

func MapAPI(cluster *v1alpha1.Cluster) Config {
	return Config{
		Name:          utils.Fallback(cluster.Spec.Title, utils.Title(cluster.Name)),
		Host:          cluster.Spec.Host,
		HTTP2:         cluster.Spec.HTTP2,
		SkipTLS:       cluster.Spec.SkipTLS,
		Certificate:   cluster.Spec.Certificate,
		SecretRef:     cluster.Spec.SecretRef,
		BasicAuth:     MapBasicAuth(cluster.Spec.BasicAuth),
		AccessControl: MapAccessControl(cluster.Spec.AccessControl),
		Plugins:       MapAPIPlugins(cluster.Spec.Plugins),
	}
}

func MapAPIPlugins(plugins []v1alpha1.Plugin) []Plugin {
	mapped := make([]Plugin, 0, len(plugins))
	for _, p := range plugins {
		mapped = append(mapped, Plugin{
			Name:        p.Source,
			Host:        p.Host,
			HTTP2:       p.HTTP2,
			SkipTLS:     p.SkipTLS,
			Certificate: p.Certificate,
			SecretRef:   p.SecretRef,
			BasicAuth:   MapBasicAuth(p.BasicAuth),
		})
	}
	return mapped
}

func MapAccessControl(ac *v1alpha1.AccessControl) AccessControl {
	if ac == nil {
		return AccessControl{}
	}

	return AccessControl{
		Emails: ac.Emails,
		Groups: ac.Groups,
	}
}

func MapBasicAuth(ba *v1alpha1.BasicAuth) BasicAuth {
	if ba == nil {
		return BasicAuth{}
	}

	return BasicAuth{
		Username: ba.Username,
		Password: ba.Password,
	}
}

func MapClusterToModel(name string, c Config) (*Cluster, error) {
	core, err := MapCoreClient(c)
	if err != nil {
		zap.L().Error("failed to create cluster client", zap.String("cluster", name), zap.Error(err))
		return nil, err
	}

	return &Cluster{
		Name:          c.Name,
		Core:          core,
		Plugins:       MapPlugins(c),
		AccessControl: c.AccessControl,
	}, nil
}

func MapCoreClient(cluster Config) (*core.Client, error) {
	options := []api.ClientOption{
		api.WithBaseURL(cluster.Host),
	}
	if cluster.HTTP2 {
		options = append(options, api.WithHTTP2())
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

	if cluster.Logging {
		options = append(options, api.WithLogging())
	}

	return core.New(options)
}

func MapPlugins(cluster Config) map[string]*plugin.Client {
	plugins := make(map[string]*plugin.Client, len(cluster.Plugins))
	for _, p := range cluster.Plugins {
		options := []api.ClientOption{
			api.WithBaseURL(p.Host),
		}

		if p.Certificate != "" {
			options = append(options, api.WithCertificate(p.Certificate))
		} else if p.SkipTLS {
			options = append(options, api.WithSkipTLS())
		}

		if p.HTTP2 {
			options = append(options, api.WithHTTP2())
		}

		if p.BasicAuth.Username != "" {
			options = append(options, api.WithBaseAuth(api.BasicAuth{
				Username: p.BasicAuth.Username,
				Password: p.BasicAuth.Password,
			}))
		}

		client, err := plugin.New(options)
		if err != nil {
			zap.L().Error("failed to create plugin client", zap.String("plugin", p.Name), zap.Error(err))
			continue
		}

		plugins[p.Name] = client
	}

	return plugins
}
