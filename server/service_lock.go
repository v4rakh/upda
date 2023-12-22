package server

type lockService interface {
	init() error
	tryLock(resource string) error
	release(resource string) error
	exists(resource string) bool
	stop()
}
