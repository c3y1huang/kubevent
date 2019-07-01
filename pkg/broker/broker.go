package broker

type BrokerOperation interface {
	Start() error
	Stop() error
	IsInitialized() bool

	Send(msg interface{}) error
}
