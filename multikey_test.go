package locker_test

import (
	"testing"

	"github.com/kamilsk/locker"
)

func TestMultiMutex(t *testing.T) {
	set := locker.MultiMutex(3)
	set.Mutex("test").Lock()
	set.Mutex("another").Lock()
	defer set.Mutex("test").Unlock()
	defer set.Mutex("another").Unlock()
	go func() { set.Mutex("test").Unlock() }()
	go func() { set.Mutex("another").Unlock() }()
	set.Mutex("test").Lock()
	set.Mutex("another").Lock()
}
