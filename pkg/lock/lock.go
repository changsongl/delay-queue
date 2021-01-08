package lock

type Locker interface {
	Lock() error
	Unlock() (bool, error)
}
