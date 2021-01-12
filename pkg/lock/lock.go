package lock

// LockerFunc locker function given a lock name,
// return a interface with Lock and Unlock method.
type LockerFunc func(name string) Locker

type Locker interface {
	Lock() error
	Unlock() (bool, error)
}
