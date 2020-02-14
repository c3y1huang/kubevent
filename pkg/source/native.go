package source

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type NativeKinds struct {
	// Type is the type of object to watch.  e.g. &v1.Pod{}
	Types []runtime.Object

	// Cache used to watch APIs
	cache cache.Cache
}

func (n NativeKinds) Start(handler handler.EventHandler, queue workqueue.RateLimitingInterface, predicates ...predicate.Predicate) error {
	if len(n.Types) == 0 {
		return fmt.Errorf("must specify NativeKinds.Types")
	}

	if n.cache == nil {
		return fmt.Errorf("must call CacheInto on NativeKinds before calling Start")
	}

	for _, obj := range n.Types {
		i, err := n.cache.GetInformer(obj)
		if err != nil {
			return err
		}

		i.AddEventHandler(eventSourceHandler{handler: handler, predicates: predicates})
	}

	return nil
}
