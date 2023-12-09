package secrets

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

type Values struct {
	Host        string `json:"api" mapstructure:"api"`
	KyvernoAPI  string `json:"kyvernoApi" mapstructure:"kyvernoApi"`
	Certificate string `json:"certificate" mapstructure:"certificate"`
	SkipTLS     bool   `json:"skipTLS" mapstructure:"skipTLS"`
	Username    string `json:"username" mapstructure:"username"`
	Password    string `json:"password" mapstructure:"password"`
}

type Client interface {
	Get(context.Context, string) (Values, error)
}

type k8sClient struct {
	client v1.SecretInterface
}

func (c *k8sClient) Get(ctx context.Context, name string) (Values, error) {
	var secret *corev1.Secret

	err := retry.OnError(retry.DefaultRetry, func(err error) bool {
		if _, ok := err.(errors.APIStatus); !ok {
			return true
		}

		if ok := errors.IsTimeout(err); ok {
			return true
		}

		if ok := errors.IsServerTimeout(err); ok {
			return true
		}

		if ok := errors.IsServiceUnavailable(err); ok {
			return true
		}

		return false
	}, func() error {
		var err error
		secret, err = c.client.Get(ctx, name, metav1.GetOptions{})

		return err
	})

	values := Values{}
	if err != nil {
		return values, err
	}

	if api, ok := secret.Data["api"]; ok {
		values.Host = string(api)
	}

	if kyvernoAPI, ok := secret.Data["kyvernoApi"]; ok {
		values.KyvernoAPI = string(kyvernoAPI)
	}

	if certificate, ok := secret.Data["certificate"]; ok {
		values.Certificate = string(certificate)
	}

	if username, ok := secret.Data["username"]; ok {
		values.Username = string(username)
	}

	if password, ok := secret.Data["password"]; ok {
		values.Password = string(password)
	}

	if skipTLS, ok := secret.Data["skipTLS"]; ok {
		v, err := strconv.ParseBool(string(skipTLS))
		if err != nil {
			zap.L().Error("failed to parse skipTLS", zap.Error(err), zap.String("secret", name))
		} else {
			values.SkipTLS = v
		}
	}

	return values, nil
}

func NewClient(secretClient v1.SecretInterface) Client {
	return &k8sClient{secretClient}
}
