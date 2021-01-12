package service

// LifeCycle service's lifecycle
type LifeCycle interface {
	Start() error
	Stop() error
}
