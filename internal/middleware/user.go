package middleware

import (
	"errors"
	"sync"

	"github.com/go-webauthn/example/internal/model"
)

func NewMemoryUserProvider() *MemoryUserProvider {
	return &MemoryUserProvider{
		users: map[string]model.User{},
		mutex: sync.RWMutex{},
	}
}

type MemoryUserProvider struct {
	users map[string]model.User
	mutex sync.RWMutex
}

func (p *MemoryUserProvider) Set(user *model.User) (err error) {
	p.mutex.Lock()

	defer p.mutex.Unlock()

	if user.Name == "" {
		return errors.New("user has no name")
	}

	p.users[user.Name] = *user

	return nil
}

func (p *MemoryUserProvider) Get(name string) (user *model.User, err error) {
	p.mutex.RLock()

	defer p.mutex.RUnlock()

	var (
		ok bool
		u  model.User
	)

	if u, ok = p.users[name]; !ok {
		return nil, errors.New("could not find user")
	}

	return &u, nil
}
