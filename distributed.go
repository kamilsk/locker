package locker

func Distributed() *dlock {
	return &dlock{}
}

type dlock struct{}

func (l *dlock) Lock(Breaker) error {
	return nil
}

func (l *dlock) Unlock(Breaker) error {
	return nil
}
