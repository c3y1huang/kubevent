package engine

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
	"sync"
)

import (
	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	batch_v1beta1 "k8s.io/api/batch/v1beta1"
	core_v1 "k8s.io/api/core/v1"
	ext_v1beta1 "k8s.io/api/extensions/v1beta1"
	storage_v1 "k8s.io/api/storage/v1"
)

type ControllerEngineFunctions interface {
	CreateController(
		name string,
		resourceTypes []string,
		eventHandlers []handler.EventHandler,
		predicater predicate.Predicate,
		reconciler reconcile.Reconciler,
	) error

	Start() error
	Stop() error
}

type ControllerEngine struct {
	Mgr       manager.Manager
	mgrMtx    sync.Mutex
	mgrStopCh chan struct{}

	apiextclient apiextclient.Interface
}

func New() (ControllerEngineFunctions, error) {
	eng := &ControllerEngine{
		mgrStopCh: make(chan struct{}),
	}

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		return nil, errors.WithMessage(err, "")
	}
	eng.Mgr = mgr
	eng.apiextclient = apiextclient.NewForConfigOrDie(mgr.GetConfig())

	return eng, nil
}

func (receiver *ControllerEngine) Start() error {
	log.Infoln("Starting controller engine serving controller management")

	if err := receiver.Mgr.Start(receiver.mgrStopCh); err != nil {
		return errors.WithMessage(err, "")
	}

	return nil
}

func (receiver *ControllerEngine) Stop() error {
	log.Infoln("Stopping controller engine")

	receiver.mgrStopCh <- struct{}{}

	return nil
}

func (receiver *ControllerEngine) CreateController(
	name string,
	resourceTypes []string,
	eventHandlers []handler.EventHandler,
	predicater predicate.Predicate,
	reconciler reconcile.Reconciler,
) error {

	receiver.injectControllerEngineAware(reconciler)
	receiver.injectControllerEngineAware(eventHandlers)

	ctrlName := fmt.Sprintf("%s-controller", name)
	log.WithField("name", ctrlName).Info("Creating controller")
	ctrl, err := controller.New(ctrlName, receiver.Mgr, controller.Options{
		Reconciler: reconciler,
	})
	if err != nil {
		return err
	}

	// v1.14: Support Workloads, Services, Config and Storage API resources implementing runtime.Object
	for _, resourceType := range resourceTypes {
		var src source.Source

		switch v := strings.ToLower(resourceType); v {
		case "pod":
			src = &source.Kind{Type: &core_v1.Pod{}}

		case "replicationcontroller":
			src = &source.Kind{Type: &core_v1.ReplicationController{}}

		case "service":
			src = &source.Kind{Type: &core_v1.Service{}}

		case "namespace":
			src = &source.Kind{Type: &core_v1.Namespace{}}

		case "persistentvolume":
			src = &source.Kind{Type: &core_v1.PersistentVolume{}}

		case "persistentvolumeclaim":
			src = &source.Kind{Type: &core_v1.PersistentVolumeClaim{}}

		case "secret":
			src = &source.Kind{Type: &core_v1.Secret{}}

		case "configmap":
			src = &source.Kind{Type: &core_v1.ConfigMap{}}

		case "endpoints":
			src = &source.Kind{Type: &core_v1.Endpoints{}}

		case "daemonset":
			src = &source.Kind{Type: &apps_v1.DaemonSet{}}

		case "statefulset":
			src = &source.Kind{Type: &apps_v1.StatefulSet{}}

		case "replicaset":
			src = &source.Kind{Type: &apps_v1.ReplicaSet{}}

		case "deployment":
			src = &source.Kind{Type: &apps_v1.Deployment{}}

		case "job":
			src = &source.Kind{Type: &batch_v1.Job{}}

		case "cronjob":
			src = &source.Kind{Type: &batch_v1beta1.CronJob{}}

		case "ingress":
			src = &source.Kind{Type: &ext_v1beta1.Ingress{}}

		case "storageclass":
			src = &source.Kind{Type: &storage_v1.StorageClass{}}

		default:
			log.WithField("kind", resourceType).Warnln("Watched resource type not supported")
			continue
		}

		for _, eventHandler := range eventHandlers {
			err := ctrl.Watch(
				src,
				eventHandler,
				predicater,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (receiver *ControllerEngine) injectControllerEngineAware(value interface{}) {
	update := func(v interface{}) {
		if aware, ok := v.(ControllerEngineAware); ok {
			aware.SetEngine(receiver)
		}
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		values := reflect.ValueOf(value)

		for i := 0; i < values.Len(); i++ {
			update(values.Index(i))
		}

	default:
		update(value)
	}
}
