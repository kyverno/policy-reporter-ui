package targetconfig

import (
	"fmt"

	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/customboard/v1alpha1"
	ui "github.com/kyverno/policy-reporter-ui/pkg/crd/client/clientset/versioned"
	informer "github.com/kyverno/policy-reporter-ui/pkg/crd/client/informers/externalversions"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
)

type Client struct {
	collection *customboard.Collection
	client     ui.Interface
}

func (c *Client) ConfigureInformer(cbInformer, ncbInformer cache.SharedIndexInformer) {
	cbInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cb := obj.(*v1alpha1.CustomBoard)
			zap.L().Info("new custom board", zap.String("name", cb.Name))

			c.collection.Add(cb.Name, customboard.MapCustomBoardToModel(cb))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			cb := newObj.(*v1alpha1.CustomBoard)
			zap.L().Info("update custom board", zap.String("name", cb.Name))

			c.collection.Add(cb.Name, customboard.MapCustomBoardToModel(cb))
		},
		DeleteFunc: func(obj interface{}) {
			cb := obj.(*v1alpha1.CustomBoard)
			zap.L().Info("delete custom board", zap.String("name", cb.Name))

			c.collection.Remove(cb.Name)
		},
	})

	ncbInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cb := obj.(*v1alpha1.NamespaceCustomBoard)
			zap.L().Info("new namespace custom board", zap.String("name", cb.Name), zap.String("namespace", cb.Namespace))

			c.collection.Add(fmt.Sprintf("%s-%s", cb.Name, cb.Namespace), customboard.MapNamespaceCustomBoardToModel(cb))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			cb := newObj.(*v1alpha1.NamespaceCustomBoard)
			zap.L().Info("update namespace custom board", zap.String("name", cb.Name), zap.String("namespace", cb.Namespace))

			c.collection.Add(fmt.Sprintf("%s-%s", cb.Name, cb.Namespace), customboard.MapNamespaceCustomBoardToModel(cb))
		},
		DeleteFunc: func(obj interface{}) {
			cb := obj.(*v1alpha1.NamespaceCustomBoard)
			zap.L().Info("delete namespace custom board", zap.String("name", cb.Name), zap.String("namespace", cb.Namespace))

			c.collection.Remove(fmt.Sprintf("%s-%s", cb.Name, cb.Namespace))
		},
	})
}

func (c *Client) Run(stopChan chan struct{}) {
	factory := informer.NewSharedInformerFactory(c.client, 0)

	cbInformer := factory.Ui().V1alpha1().CustomBoards().Informer()
	ncbInformer := factory.Ui().V1alpha1().NamespaceCustomBoards().Informer()

	c.ConfigureInformer(cbInformer, ncbInformer)

	factory.Start(stopChan)

	if !cache.WaitForCacheSync(stopChan, cbInformer.HasSynced) {
		zap.L().Error("Failed to sync custom board cache")
		return
	}

	if !cache.WaitForCacheSync(stopChan, ncbInformer.HasSynced) {
		zap.L().Error("Failed to sync namespace custom board cache")
		return
	}

	zap.L().Info("custom board cache synced")
}

func NewClient(client ui.Interface, targets *customboard.Collection) *Client {

	return &Client{
		client:     client,
		collection: targets,
	}
}
