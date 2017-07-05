package cache

import (
	"github.com/stretchr/testify/mock"
)

type MockCredentialsCache struct {
	mock.Mock

	SetFn func(registry string, entry *AuthEntry)
}

func (c *MockCredentialsCache) Get(registry string) *AuthEntry {
	args := c.Called(registry)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(*AuthEntry)
}

func (c *MockCredentialsCache) Set(registry string, entry *AuthEntry) {
	if c.SetFn != nil {
		c.SetFn(registry, entry)
		return
	}
	c.Called(registry, entry)
}

func (c *MockCredentialsCache) List() []*AuthEntry {
	args := c.Called()
	return args.Get(0).([]*AuthEntry)
}

func (c *MockCredentialsCache) Clear() {
	c.Called()
}
