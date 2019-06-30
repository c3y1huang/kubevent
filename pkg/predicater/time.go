package predicater

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"
)

type Time struct {
	startTime time.Time
}

func NewTime(t time.Time) predicate.Predicate {
	return &Time{
		startTime: t,
	}
}

func (receiver *Time) Create(e event.CreateEvent) bool {
	return e.Meta.GetCreationTimestamp().After(receiver.startTime)
}

func (receiver *Time) Delete(e event.DeleteEvent) bool {
	return e.Meta.GetDeletionTimestamp().After(receiver.startTime)
}

func (receiver *Time) Update(e event.UpdateEvent) bool {
	return e.MetaNew.GetCreationTimestamp().After(receiver.startTime)
}

func (receiver *Time) Generic(e event.GenericEvent) bool {
	return e.Meta.GetCreationTimestamp().After(receiver.startTime)
}
