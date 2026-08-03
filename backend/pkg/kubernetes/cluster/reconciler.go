package cluster

import (
	"context"

	"github.com/kyverno/policy-reporter-ui/pkg/cluster"
	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/ui/v1alpha1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type reconciler struct {
	client     client.Client
	collection *cluster.Collection
	loader     *cluster.SecretLoader
}

func newReconciler(client client.Client, collection *cluster.Collection, loader *cluster.SecretLoader) *reconciler {
	return &reconciler{
		client:     client,
		collection: collection,
		loader:     loader,
	}
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error

	cl := &v1alpha1.Cluster{}
	err = r.client.Get(ctx, req.NamespacedName, cl)
	if err == nil {
		config := cluster.MapAPI(cl)
		config = r.loader.LoadConfig(ctx, config)
		zap.L().Debug("cluster added", zap.String("name", cl.Name))
		r.collection.Add(cl.Name, config)
	} else if errors.IsNotFound(err) {
		zap.L().Debug("cluster removed", zap.String("name", req.NamespacedName.Name))
		r.collection.Remove(req.NamespacedName.Name)
		return ctrl.Result{}, nil
	} else {
		zap.L().Error("failed to get cluster", zap.String("name", req.NamespacedName.Name), zap.Error(err))
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
