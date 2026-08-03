package customboard

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/ui/v1alpha1"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
)

type Client struct {
	collection *customboard.Collection
	mgr        manager.Manager
}

func (c *Client) Start(ctx context.Context) error {
	cbBuilder := ctrl.NewControllerManagedBy(c.mgr).For(&v1alpha1.CustomBoard{})
	ncbBuilder := ctrl.NewControllerManagedBy(c.mgr).For(&v1alpha1.NamespaceCustomBoard{})

	if err := cbBuilder.Complete(newReconciler(c.mgr.GetClient(), c.collection)); err != nil {
		return fmt.Errorf("failed to construct custom board manager: %w", err)
	}
	if err := ncbBuilder.Complete(newReconciler(c.mgr.GetClient(), c.collection)); err != nil {
		return fmt.Errorf("failed to construct namespace custom board manager: %w", err)
	}
	return nil
}

func NewClient(mgr manager.Manager, targets *customboard.Collection) (*Client, error) {
	return &Client{
		collection: targets,
		mgr:        mgr,
	}, nil
}
