package locker

func Semaphore(capacity int) *semaphore {
	return &semaphore{}
}

type semaphore struct{}

func (s *semaphore) Lock(Breaker) error {
	return nil
}

func (s *semaphore) Unlock(Breaker) error {
	return nil
}
