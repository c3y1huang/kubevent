package predicater

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

type TimePredicater struct {
	startTime time.Time
}

func NewTimePredicater(t time.Time) predicate.Predicate {
	return &TimePredicater{
		startTime: t,
	}
}

func (t *TimePredicater) Create(e event.CreateEvent) bool {
	return e.Meta.GetCreationTimestamp().After(t.startTime)
}

func (t *TimePredicater) Delete(e event.DeleteEvent) bool {
	return e.Meta.GetDeletionTimestamp().After(t.startTime)
}

func (t *TimePredicater) Update(e event.UpdateEvent) bool {
	return e.MetaNew.GetCreationTimestamp().After(t.startTime)
}

func (t *TimePredicater) Generic(e event.GenericEvent) bool {
	return e.Meta.GetCreationTimestamp().After(t.startTime)
}
