package source

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	toolscache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/event"
	controllerruntimehandler "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// eventSourceHandler uses the same logic from sigs.k8s.io/controller-runtime@v0.4.0/pkg/source/internal/eventsource.go but without queue
type eventSourceHandler struct {
	handler    controllerruntimehandler.EventHandler
	predicates []predicate.Predicate
}

func (e eventSourceHandler) OnAdd(obj interface{}) {
	c := event.CreateEvent{}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(obj); err == nil {
		c.Meta = o
	} else {
		logrus.WithError(err).Errorln("OnAdd missing Meta", "object", obj, "type", fmt.Sprintf("%T", obj))
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := obj.(runtime.Object); ok {
		c.Object = o
	} else {
		logrus.Errorln("OnAdd missing runtime.Object", "object", obj, "type", fmt.Sprintf("%T", obj))
		return
	}

	for _, p := range e.predicates {
		if !p.Create(c) {
			return
		}
	}

	// Invoke create handler
	e.handler.Create(c, nil)
}

func (e eventSourceHandler) OnUpdate(oldObj, newObj interface{}) {
	u := event.UpdateEvent{}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(oldObj); err == nil {
		u.MetaOld = o
	} else {
		logrus.WithError(err).Errorln("OnUpdate missing MetaOld", "object", oldObj, "type", fmt.Sprintf("%T", oldObj))
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := oldObj.(runtime.Object); ok {
		u.ObjectOld = o
	} else {
		logrus.Errorln("OnUpdate missing ObjectOld", "object", oldObj, "type", fmt.Sprintf("%T", oldObj))
		return
	}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(newObj); err == nil {
		u.MetaNew = o
	} else {
		logrus.WithError(err).Errorln("OnUpdate missing MetaNew", "object", newObj, "type", fmt.Sprintf("%T", newObj))
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := newObj.(runtime.Object); ok {
		u.ObjectNew = o
	} else {
		logrus.Errorln("OnUpdate missing ObjectNew", "object", oldObj, "type", fmt.Sprintf("%T", oldObj))
		return
	}

	for _, p := range e.predicates {
		if !p.Update(u) {
			return
		}
	}

	// Invoke update handler
	e.handler.Update(u, nil)
}

func (e eventSourceHandler) OnDelete(obj interface{}) {
	d := event.DeleteEvent{}

	// Deal with tombstone events by pulling the object out.  Tombstone events wrap the object in a
	// DeleteFinalStateUnknown struct, so the object needs to be pulled out.
	// Copied from sample-controller
	// This should never happen if we aren't missing events, which we have concluded that we are not
	// and made decisions off of this belief.  Maybe this shouldn't be here?
	var ok bool
	if _, ok = obj.(metav1.Object); !ok {
		// If the object doesn't have Metadata, assume it is a tombstone object of type DeletedFinalStateUnknown
		tombstone, ok := obj.(toolscache.DeletedFinalStateUnknown)
		if !ok {
			logrus.Errorln("Error decoding objects.  Expected Cache.DeletedFinalStateUnknown",
				"type", fmt.Sprintf("%T", obj),
				"object", obj)
			return
		}

		// Set obj to the tombstone obj
		obj = tombstone.Obj
	}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(obj); err == nil {
		d.Meta = o
	} else {
		logrus.WithError(err).Errorln("OnDelete missing Meta", "object", obj, "type", fmt.Sprintf("%T", obj))
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := obj.(runtime.Object); ok {
		d.Object = o
	} else {
		logrus.Errorln("OnDelete missing runtime.Object", "object", obj, "type", fmt.Sprintf("%T", obj))
		return
	}

	for _, p := range e.predicates {
		if !p.Delete(d) {
			return
		}
	}

	// Invoke delete handler
	e.handler.Delete(d, nil)
}
