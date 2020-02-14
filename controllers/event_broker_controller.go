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
	"k8s.io/apimachinery/pkg/runtime/schema"
	controllerruntimecontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	controllerruntimehandler "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeventv1alpha1 "github.com/innobead/kubevent/api/v1alpha1"
)

var passReconciler = reconciler.NewPassReconciler()

// EventBrokerController reconciles a EventBroker object
type EventBrokerController struct {
	client.Client
	Scheme           *runtime.Scheme
	Mgr              manager.Manager
	watchControllers map[string]eventBrokerWatchController
}

type eventBrokerWatchController struct {
	stop chan struct{}
	// controller is runtime controller
	controller           controllerruntimecontroller.Controller
	eventBrokerOperation handler.Operation
}

func NewEventBrokerController(manager manager.Manager) (*EventBrokerController, error) {
	controller := &EventBrokerController{
		Client:           manager.GetClient(),
		Scheme:           manager.GetScheme(),
		Mgr:              manager,
		watchControllers: map[string]eventBrokerWatchController{},
	}

	err := controller.SetupWithManager(manager)

	return controller, err
}

// +kubebuilder:rbac:groups=kubevent.innobead,resources=eventbroker,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kubevent.innobead,resources=eventbroker/status,verbs=get;update;patch

func (e *EventBrokerController) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	logger := logrus.WithField("request", req).Logger

	logger.Infoln("reconciling request")

	ctx := context.Background()
	eventBroker := &kubeventv1alpha1.EventBroker{}

	// delete eventBrokerWatchController if not found
	if err := e.Client.Get(ctx, req.NamespacedName, eventBroker); err != nil {
		logger.WithError(err).Errorln("failed to get object of event broker watch controller")

		e.deleteEventBrokerWatchController(logger, eventBroker)

		return ctrl.Result{}, nil
	}

	// delete eventBrokerWatchController if found, then create a new one for update case
	e.deleteEventBrokerWatchController(logger, eventBroker)

	// create eventBrokerWatchController
	logger.Infoln("creating event broker watch controller")

	controllerName := getEventBrokerName(eventBroker)
	controllerLogger := logger.WithField("controller", controllerName).Logger
	controller, err := controllerruntimecontroller.New(
		controllerName,
		e.Mgr,
		controllerruntimecontroller.Options{
			MaxConcurrentReconciles: 1,
			Reconciler:              passReconciler,
		},
	)
	if err != nil {
		controllerLogger.Errorln("failed to create event broker watch controller")
		return createErrorResult(), err
	}

	stopChan := make(chan struct{})
	e.watchControllers[controllerName] = eventBrokerWatchController{
		stop:       stopChan,
		controller: controller,
	}

	go func() {
		controllerLogger.Infoln("starting event broker watch controller")

		if err := controller.Start(stopChan); err != nil {
			controllerLogger.WithError(err).Warnln("failed to start event broker watch controller")
		}
	}()

	controllerLogger.Infoln("watching event")

	watchController := e.watchControllers[controllerName]
	var eventBrokerHandler controllerruntimehandler.EventHandler

	eventBrokerHandler, watchController.eventBrokerOperation, err = handler.CreateEventBrokerHandler(&eventBroker.Spec)
	if err != nil {
		controllerLogger.WithError(err).Errorln("failed to start kafka event eventBroker event handler")
		return ctrl.Result{}, err
	}

	go func() {
		if err := watchController.eventBrokerOperation.Start(); err != nil {
			controllerLogger.WithError(err).Errorln("failed to start kafka event eventBroker event handler")
		}
	}()

	var gvks []schema.GroupVersionKind
	if eventBroker.Spec.WatchAllResources {
		for k := range e.Scheme.AllKnownTypes() {
			gvks = append(gvks, k)
		}
	} else {
		gvks = append(gvks, util.ToSchemaGroupVersionKinds(eventBroker.Spec.WatchResources)...)
	}

	if err = controller.Watch(
		&source.DynamicKinds{GroupVersionKinds: gvks, Cache: e.Mgr.GetCache()},
		eventBrokerHandler,
		predicater.NewTimePredicater(time.Now()),
	); err != nil {
		controllerLogger.WithError(err).Errorln("failed to watch event")

		return createErrorResult(), err
	}

	return ctrl.Result{}, nil
}

func (e *EventBrokerController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeventv1alpha1.EventBroker{}).
		Complete(e)
}

func (e *EventBrokerController) deleteEventBrokerWatchController(logger *logrus.Logger, broker *kubeventv1alpha1.EventBroker) {
	name := getEventBrokerName(broker)

	logger.Infoln("deleting event broker watch controller")

	if controller, ok := e.watchControllers[name]; ok {
		defer func() {
			close(controller.stop)
		}()

		defer func() {
			if err := controller.eventBrokerOperation.Stop(); err != nil {
				logger.WithError(err).Errorln("failed to stop kafka event broker event handler")
			}
		}()

		controller.stop <- struct{}{}

		delete(e.watchControllers, name)
	}
}

func createErrorResult() ctrl.Result {
	return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}
}

func getEventBrokerName(broker *kubeventv1alpha1.EventBroker) string {
	return fmt.Sprintf("event-broker-%s", broker.Name)
}
