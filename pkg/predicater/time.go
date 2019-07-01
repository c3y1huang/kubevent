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

func (receiver *TimePredicater) Create(e event.CreateEvent) bool {
	return e.Meta.GetCreationTimestamp().After(receiver.startTime)
}

func (receiver *TimePredicater) Delete(e event.DeleteEvent) bool {
	return e.Meta.GetDeletionTimestamp().After(receiver.startTime)
}

func (receiver *TimePredicater) Update(e event.UpdateEvent) bool {
	return e.MetaNew.GetCreationTimestamp().After(receiver.startTime)
}

func (receiver *TimePredicater) Generic(e event.GenericEvent) bool {
	return e.Meta.GetCreationTimestamp().After(receiver.startTime)
}
