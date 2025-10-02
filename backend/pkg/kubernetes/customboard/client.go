package targetconfig

import (
	"fmt"

	"github.com/kyverno/policy-reporter-ui/pkg/crd/api/customboard/v1alpha1"
	ui "github.com/kyverno/policy-reporter-ui/pkg/crd/client/clientset/versioned"
	cb "github.com/kyverno/policy-reporter-ui/pkg/crd/client/clientset/versioned/typed/customboard/v1alpha1"
	informer "github.com/kyverno/policy-reporter-ui/pkg/crd/client/informers/externalversions"
	"github.com/kyverno/policy-reporter-ui/pkg/customboard"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
)

type Client struct {
	collection  *customboard.Collection
	cbInformer  cache.SharedIndexInformer
	ncbInformer cache.SharedIndexInformer
	client      cb.CustomBoardInterface
}

func (c *Client) ConfigureInformer() {
	c.cbInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			tc := obj.(*v1alpha1.CustomBoard)
			zap.L().Info("new custom board", zap.String("name", tc.Name))

			c.collection.Add(tc.Name, customboard.MapCustomBoardToModel(tc))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			tc := newObj.(*v1alpha1.CustomBoard)
			zap.L().Info("update custom board", zap.String("name", tc.Name))

			c.collection.Add(tc.Name, customboard.MapCustomBoardToModel(tc))
		},
		DeleteFunc: func(obj interface{}) {
			tc := obj.(*v1alpha1.CustomBoard)
			zap.L().Info("delete custom board", zap.String("name", tc.Name))

			c.collection.Remove(tc.Name)
		},
	})

	c.ncbInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			tc := obj.(*v1alpha1.NamespaceCustomBoard)
			zap.L().Info("new namespace custom board", zap.String("name", tc.Name))

			c.collection.Add(fmt.Sprintf("%s/%s", tc.Name, tc.Namespace), customboard.MapNamespaceCustomBoardToModel(tc))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			tc := newObj.(*v1alpha1.NamespaceCustomBoard)
			zap.L().Info("update namespace custom board", zap.String("name", tc.Name))

			c.collection.Add(fmt.Sprintf("%s/%s", tc.Name, tc.Namespace), customboard.MapNamespaceCustomBoardToModel(tc))
		},
		DeleteFunc: func(obj interface{}) {
			tc := obj.(*v1alpha1.NamespaceCustomBoard)
			zap.L().Info("delete namespace custom board", zap.String("name", tc.Name))

			c.collection.Remove(fmt.Sprintf("%s/%s", tc.Name, tc.Namespace))
		},
	})
}

func (c *Client) Run(stopChan chan struct{}) {
	go c.cbInformer.Run(stopChan)

	if !cache.WaitForCacheSync(stopChan, c.cbInformer.HasSynced) {
		zap.L().Error("Failed to sync custom board cache")
		return
	}

	go c.ncbInformer.Run(stopChan)

	if !cache.WaitForCacheSync(stopChan, c.ncbInformer.HasSynced) {
		zap.L().Error("Failed to sync namespace custom board cache")
		return
	}

	zap.L().Info("custom board cache synced")
}

func NewClient(client ui.Interface, targets *customboard.Collection) *Client {
	factory := informer.NewSharedInformerFactory(client, 0)

	return &Client{
		cbInformer:  factory.Ui().V1alpha1().CustomBoards().Informer(),
		ncbInformer: factory.Ui().V1alpha1().NamespaceCustomBoards().Informer(),
		collection:  targets,
	}
}
