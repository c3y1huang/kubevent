package source

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type DynamicKinds struct {
	// Type is the type of object to watch.  e.g. &v1.Pod{}
	GroupVersionKinds []schema.GroupVersionKind

	// cache used to watch APIs
	cache cache.Cache
}

func (d DynamicKinds) Start(handler handler.EventHandler, queue workqueue.RateLimitingInterface, predicates ...predicate.Predicate) error {
	if len(d.GroupVersionKinds) == 0 {
		return fmt.Errorf("must specify DynamicKinds.GroupVersionKinds")
	}

	if d.cache == nil {
		return fmt.Errorf("must call CacheInto on NativeKinds before calling Start")
	}

	for _, kind := range d.GroupVersionKinds {
		i, err := d.cache.GetInformerForKind(kind)
		if err != nil {
			return err
		}

		i.AddEventHandler(eventSourceHandler{handler: handler, predicates: predicates})
	}

	return nil
}
