/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/innobead/kubevent/pkg/handler"
	"github.com/innobead/kubevent/pkg/predicater"
	"github.com/innobead/kubevent/pkg/reconciler"
	"github.com/innobead/kubevent/pkg/source"
	"github.com/innobead/kubevent/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeventv1alpha1 "github.com/innobead/kubevent/api/v1alpha1"
)

// BrokerReconciler reconciles a Broker object
type BrokerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Mgr    manager.Manager
}

// +kubebuilder:rbac:groups=kubevent.innobead,resources=brokers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubevent.innobead,resources=brokers/status,verbs=get;update;patch

func (r *BrokerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	broker := &kubeventv1alpha1.Broker{}
	if err := r.Client.Get(ctx, req.NamespacedName, broker); err != nil {
		return createErrorResult(), err
	}

	name := fmt.Sprintf("event-%s-broker-%s", broker.Spec.Provider, broker.Name)
	c, err := controller.New(name, r.Mgr, controller.Options{MaxConcurrentReconciles: 1, Reconciler: reconciler.NewPassReconciler()})
	if err != nil {
		return createErrorResult(), err
	}

	gvks := util.ToSchemaGroupVersionKinds(broker.Spec.GroupVersionKinds)
	if err = c.Watch(&source.DynamicKinds{GroupVersionKinds: gvks}, &handler.KafkaHandler{}, predicater.NewTimePredicater(time.Now())); err != nil {
		return createErrorResult(), err
	}

	return ctrl.Result{}, nil
}

func (r *BrokerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeventv1alpha1.Broker{}).
		Complete(r)
}

func createErrorResult() ctrl.Result {
	return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}
}
