package customboard

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/ui/v1alpha1"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
)

type reconciler struct {
	client     client.Client
	collection *customboard.Collection
}

func newReconciler(client client.Client, collection *customboard.Collection) *reconciler {
	return &reconciler{
		client:     client,
		collection: collection,
	}
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error

	if req.Namespace != "" {
		ncb := &v1alpha1.NamespaceCustomBoard{}
		err = r.client.Get(ctx, req.NamespacedName, ncb)
		if err == nil {
			r.collection.Add(fmt.Sprintf("%s-%s", ncb.Name, ncb.Namespace), customboard.MapNamespaceCustomBoardToModel(ncb))
		} else if errors.IsNotFound(err) {
			r.collection.Remove(fmt.Sprintf("%s-%s", ncb.Name, ncb.Namespace))
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	} else {
		cb := &v1alpha1.CustomBoard{}
		err = r.client.Get(ctx, req.NamespacedName, cb)
		if err == nil {
			r.collection.Add(cb.Name, customboard.MapCustomBoardToModel(cb))
		} else if errors.IsNotFound(err) {
			r.collection.Remove(cb.Name)
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
