package reconciler

import (
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DummyReconciler struct {
}

func (*DummyReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log.Debugf("Reconciling %v", req)

	return reconcile.Result{}, nil
}
