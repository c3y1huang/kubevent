package engine

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
	"sync"
)

type ControllerEngineFunctions interface {
	CreateController(
		watchedApiTypes []string,
		eventHandlers []handler.EventHandler,
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
	if err := receiver.Mgr.Start(receiver.mgrStopCh); err != nil {
		return errors.WithMessage(err, "")
	}

	<-receiver.mgrStopCh

	return nil
}

func (receiver *ControllerEngine) Stop() error {
	receiver.mgrStopCh <- struct{}{}

	return nil
}

func (receiver *ControllerEngine) CreateController(
	watchedApiTypes []string,
	eventHandlers []handler.EventHandler,
	reconciler reconcile.Reconciler,
) error {

	receiver.injectControllerEngineAware(reconciler)
	receiver.injectControllerEngineAware(eventHandlers)

	bu := builder.ControllerManagedBy(receiver.Mgr)
	for _, watchedApiType := range watchedApiTypes {
		var s source.Source

		switch v := strings.ToLower(watchedApiType); v {
		case "pod":
			s = &source.Kind{Type: &v1.Pod{}}
		}

		for _, eventHandler := range eventHandlers {
			bu = bu.Watches(
				s,
				eventHandler,
			)
		}
	}

	if err := bu.Complete(reconciler); err != nil {
		return errors.WithMessage(err, "")
	}

	if receiver.mgrStopCh != nil {
		receiver.mgrMtx.Lock()
		defer receiver.mgrMtx.Unlock()

		if err := receiver.Stop(); err != nil {
			return err
		}

		if err := receiver.Start(); err != nil {
			return err
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
		for _, v := range value.([]interface{}) {
			update(v)
		}

	default:
		update(value)
	}
}
