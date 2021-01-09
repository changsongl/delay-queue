package lock

type LockerFunc func(name string) Locker

type Locker interface {
	Lock() error
	Unlock() (bool, error)
}
