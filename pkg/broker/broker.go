package broker

type Operation interface {
	Start() error
	Stop() error
	Send(msg interface{}) error
}
