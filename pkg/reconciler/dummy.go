package reconciler

import (
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DummyReconciler struct {
}

func NewDummy() reconcile.Reconciler {
	return &DummyReconciler{}
}

func (*DummyReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log.WithField("request", req).Info("Reconciling") // TODO change to Debug

	return reconcile.Result{}, nil
}
