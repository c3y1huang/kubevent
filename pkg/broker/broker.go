package broker

type Operation interface {
	Start() error
	Stop() error
	IsInitialized() bool

	Send(msg interface{}) error
}
