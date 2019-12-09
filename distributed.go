package locker

func Distributed() *distributed {
	return &distributed{}
}

type distributed struct{}

func (d *distributed) Lock(Breaker) error {
	return nil
}

func (d *distributed) Unlock(Breaker) error {
	return nil
}
