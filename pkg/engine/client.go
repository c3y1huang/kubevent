package engine

import "sigs.k8s.io/controller-runtime/pkg/client"

type ControllerEngineAware interface {
	SetEngine(engine *ControllerEngine)
}

type ControllerEngineAwareType struct {
	Eng    *ControllerEngine
	Client client.Client
}

func (rec *ControllerEngineAwareType) SetEngine(engine *ControllerEngine) {
	rec.Eng = engine
	rec.Client = engine.Mgr.GetClient()
}
