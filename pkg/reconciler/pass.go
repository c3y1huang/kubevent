package reconciler

import (
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PassReconciler struct {
}

func NewPassReconciler() reconcile.Reconciler {
	return &PassReconciler{}
}

func (*PassReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log.WithField("request", req).Debugln("reconciled")

	return reconcile.Result{}, nil
}
