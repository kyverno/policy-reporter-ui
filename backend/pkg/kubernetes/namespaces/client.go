package namespaces

import (
	"context"

	"github.com/kyverno/policy-reporter-ui/pkg/kubernetes"
	"github.com/kyverno/policy-reporter-ui/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client interface {
	List(context.Context, map[string]string) ([]string, error)
}

type k8sClient struct {
	client v1.NamespaceInterface
}

func (c *k8sClient) List(ctx context.Context, selector map[string]string) ([]string, error) {
	return kubernetes.Retry(func() ([]string, error) {
		namespaces, err := c.client.List(ctx, metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: selector})})
		if err != nil {
			return nil, err
		}

		return utils.Map(namespaces.Items, func(ns corev1.Namespace) string {
			return ns.Name
		}), nil
	})
}

func NewClient(secretClient v1.NamespaceInterface) Client {
	return &k8sClient{secretClient}
}
