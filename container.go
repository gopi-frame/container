// Package container provides an implementation of [github.com/gopi-frame/contract/container.Container]
package container

import (
	"fmt"
	"github.com/gopi-frame/collection/kv"
	"github.com/gopi-frame/exception"
)

// Container container
type Container[T any] struct {
	instances    *kv.Map[string, T]
	constructors *kv.Map[string, func() (T, error)]
}

// New creates a new container
func New[T any]() *Container[T] {
	return &Container[T]{
		instances:    kv.NewMap[string, T](),
		constructors: kv.NewMap[string, func() (T, error)](),
	}
}

// Set sets a value in the container by name
func (c *Container[T]) Set(name string, value T) {
	c.instances.Lock()
	defer c.instances.Unlock()
	c.instances.Set(name, value)
}

// Get gets a value from the container by name
func (c *Container[T]) Get(name string) T {
	if v, err := c.GetE(name); err == nil {
		return v
	}
	return *(new(T))
}

// GetE gets a value from the container by name or error if not found
func (c *Container[T]) GetE(name string) (T, error) {
	locked := c.instances.TryRLock()
	v, ok := c.instances.Get(name)
	if ok {
		if locked {
			c.instances.RUnlock()
		}
		return v, nil
	}
	if locked {
		c.instances.RUnlock()
	}
	v, err := c.Make(name)
	if err != nil {
		return *(new(T)), err
	}
	defer func() {
		if c.instances.TryLock() {
			defer c.instances.Unlock()
		}
		c.instances.Set(name, v)
	}()
	return v, nil
}

// Has checks if the container has a value or a constructor by name
func (c *Container[T]) Has(name string) bool {
	if c.instances.TryRLock() {
		defer c.instances.RUnlock()
	}
	ok := c.instances.ContainsKey(name)
	if ok {
		return true
	}
	if c.constructors.TryRLock() {
		defer c.constructors.RUnlock()
	}
	return c.constructors.ContainsKey(name)
}

// Lazy sets a constructor function that will be called when the value is accessed
func (c *Container[T]) Lazy(name string, constructor func() (T, error)) {
	c.constructors.Lock()
	defer c.constructors.Unlock()
	c.constructors.Set(name, constructor)
}

// Make makes a new instance of the value.
func (c *Container[T]) Make(name string) (T, error) {
	if c.constructors.TryRLock() {
		defer c.constructors.RUnlock()
	}
	constructor, ok := c.constructors.Get(name)
	if ok {
		return constructor()
	}
	return *(new(T)),
		exception.NewArgumentException("name", name, fmt.Sprintf("constructor for %s not found", name))
}
