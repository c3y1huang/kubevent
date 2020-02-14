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
	"github.com/sirupsen/logrus"
	controllerruntimecontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeventv1alpha1 "github.com/innobead/kubevent/api/v1alpha1"
)

// EventBrokerController reconciles a EventBroker object
type EventBrokerController struct {
	client.Client
	Scheme      *runtime.Scheme
	Mgr         manager.Manager
	controllers map[string]eventBrokerWatchController
}

type eventBrokerWatchController struct {
	stop                 chan struct{}
	controller           controllerruntimecontroller.Controller
	eventBrokerOperation handler.Operation
	A                    string
}

func NewEventBrokerController(manager manager.Manager) (*EventBrokerController, error) {
	controller := &EventBrokerController{
		Client:      manager.GetClient(),
		Scheme:      manager.GetScheme(),
		Mgr:         manager,
		controllers: map[string]eventBrokerWatchController{},
	}

	err := controller.SetupWithManager(manager)

	return controller, err
}

// +kubebuilder:rbac:groups=kubevent.innobead,resources=eventbroker,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubevent.innobead,resources=eventbroker/status,verbs=get;update;patch

func (r *EventBrokerController) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	logrus.WithField("request", req).Infoln("reconciling request")

	ctx := context.Background()
	broker := &kubeventv1alpha1.EventBroker{}

	// delete eventBrokerWatchController if not found
	if err := r.Client.Get(ctx, req.NamespacedName, broker); err != nil {
		//TODO delete controller if object not found
		logrus.WithError(err).WithField("request", req).Errorln("failed to get object of event broker watch controller")

		r.deleteEventBrokerWatchController(broker)

		return ctrl.Result{}, nil
	}

	// delete eventBrokerWatchController if found, then create a new one for update case
	r.deleteEventBrokerWatchController(broker)

	// create eventBrokerWatchController
	logrus.WithField("req", req).Infoln("creating event broker watch controller")

	brokerName := getEventBrokerName(broker)
	controller, err := controllerruntimecontroller.New(
		brokerName,
		r.Mgr,
		controllerruntimecontroller.Options{
			MaxConcurrentReconciles: 1,
			Reconciler:              reconciler.NewPassReconciler(),
		},
	)
	if err != nil {
		logrus.WithField("req", req).Errorln("failed to create event broker watch controller")
		return createErrorResult(), err
	}

	stopChan := make(chan struct{})
	r.controllers[brokerName] = eventBrokerWatchController{
		stop:       stopChan,
		controller: controller,
	}

	go func() {
		if err := controller.Start(stopChan); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"req":    req,
				"broker": brokerName,
			}).Errorln("failed to start event broker watch controller")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"req":    req,
		"broker": brokerName,
	}).Infoln("watching event")

	if broker.Spec.Kafka != nil {
		eventBrokerHandler := handler.NewKafkaHandler(broker.Spec.Kafka)
		c := r.controllers[brokerName]
		c.eventBrokerOperation = handler.Operation(eventBrokerHandler)

		go func() {
			if err := eventBrokerHandler.Start(); err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"req":    req,
					"broker": brokerName,
				}).Errorln("failed to start kafka event broker event handler")
			}
		}()

		gvks := util.ToSchemaGroupVersionKinds(broker.Spec.GroupVersionKinds)
		if err = controller.Watch(
			&source.DynamicKinds{GroupVersionKinds: gvks, Cache: r.Mgr.GetCache()},
			eventBrokerHandler,
			predicater.NewTimePredicater(time.Now()),
		); err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"name":   req.NamespacedName,
				"broker": brokerName,
			}).Errorln("failed to watch event")

			return createErrorResult(), err
		}
	}

	return ctrl.Result{}, nil
}

func (r *EventBrokerController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeventv1alpha1.EventBroker{}).
		Complete(r)
}

func (r *EventBrokerController) deleteEventBrokerWatchController(broker *kubeventv1alpha1.EventBroker) {
	name := getEventBrokerName(broker)

	logrus.WithField("name", name).Infoln("deleting event broker watch controller")

	if controller, ok := r.controllers[name]; ok {
		defer func() {
			close(controller.stop)
		}()

		defer func() {
			if err := controller.eventBrokerOperation.Stop(); err != nil {
				logrus.WithError(err).WithField("name", name).Errorln("failed to stop kafka event broker event handler")
			}
		}()

		controller.stop <- struct{}{}

		delete(r.controllers, name)
	}
}

func createErrorResult() ctrl.Result {
	return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}
}

func getEventBrokerName(broker *kubeventv1alpha1.EventBroker) string {
	return fmt.Sprintf("event-broker-%s", broker.Name)
}
