package cluster

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kyverno/policy-reporter-ui/pkg/cluster"
	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/ui/v1alpha1"
)

type Client struct {
	collection *cluster.Collection
	mgr        manager.Manager
	loader     *cluster.SecretLoader
}

func (c *Client) Start(ctx context.Context) error {
	builder := ctrl.NewControllerManagedBy(c.mgr).For(&v1alpha1.Cluster{})

	if err := builder.Complete(newReconciler(c.mgr.GetClient(), c.collection, c.loader)); err != nil {
		return fmt.Errorf("failed to construct cluster manager: %w", err)
	}

	return nil
}

func NewClient(mgr manager.Manager, targets *cluster.Collection, loader *cluster.SecretLoader) (*Client, error) {
	return &Client{
		collection: targets,
		mgr:        mgr,
		loader:     loader,
	}, nil
}
