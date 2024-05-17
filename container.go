package container

import (
	"sync"

	"github.com/gopi-frame/contract/container"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/support/lists"
	"github.com/gopi-frame/support/maps"
)

var _ container.Container = (*Container)(nil)

// NewContainer new container
func NewContainer() *Container {
	container := &Container{
		instances: maps.NewMap[string, any](),
		aliases:   maps.NewMap[string, string](),
		bindings:  maps.NewMap[string, func(container.Container) any](),
		extends:   maps.NewMap[string, *lists.List[func(instance any, container container.Container) any]](),
	}
	return container
}

// Container container interface
type Container struct {
	sync.Mutex
	instances *maps.Map[string, any]
	aliases   *maps.Map[string, string]
	bindings  *maps.Map[string, func(container.Container) any]
	extends   *maps.Map[string, *lists.List[func(instance any, container container.Container) any]]
}

func (c *Container) isAlias(name string) bool {
	return c.aliases.ContainsKey(name)
}

func (c *Container) realName(name string) string {
	if alias, ok := c.aliases.Get(name); ok {
		return c.realName(alias)
	}
	return name
}

func (c *Container) resolve(abstract string) any {
	abstract = c.realName(abstract)
	if instance, ok := c.instances.Get(abstract); ok {
		return instance
	}
	concrate, ok := c.bindings.Get(abstract)
	if !ok {
		panic("concrate not found")
	}
	instance := concrate(c)
	if extenders, ok := c.extends.Get(abstract); ok && extenders != nil {
		extenders.Each(func(index int, extender func(instance any, container container.Container) any) bool {
			instance = extender(instance, c)
			return true
		})
	}
	c.Set(abstract, instance)
	return instance
}

func (c *Container) has(abstract string) bool {
	return c.isAlias(abstract) || c.instances.ContainsKey(abstract) || c.bindings.ContainsKey(abstract)
}

// Has has
func (c *Container) Has(abstract string) bool {
	if c.TryLock() {
		defer c.Unlock()
	}
	return c.has(abstract)
}

// Get get an instance
func (c *Container) Get(abstract string) any {
	if c.TryLock() {
		defer c.Unlock()
	}
	var instance any
	exception.Try(func() {
		instance = c.resolve(abstract)
	}).CatchAll(func(err error) {
		if c.has(abstract) {
			panic(err)
		}
		panic(NewEntryNotFoundException(abstract))
	}).Run()
	return instance
}

// Bind bind
func (c *Container) Bind(abstract string, concrate func(container.Container) any) {
	if c.TryLock() {
		defer c.Unlock()
	}
	c.remove(abstract)
	c.removeAlias(abstract)
	if concrate == nil {
		panic(exception.NewEmptyArgumentException("concrate"))
	}
	c.bindings.Set(abstract, concrate)
}

// BindIf bind if abstract is not bound
func (c *Container) BindIf(abstract string, concrate func(container.Container) any) {
	if c.TryLock() {
		defer c.Unlock()
	}
	if c.bindings.ContainsKey(abstract) {
		return
	}
	c.Bind(abstract, concrate)
}

func (c *Container) set(abstract string, instance any) {
	c.remove(abstract)
	c.removeAlias(abstract)
	c.instances.Set(abstract, instance)
}

// Set register an existing instance
func (c *Container) Set(abstract string, instance any) {
	if c.TryLock() {
		defer c.Unlock()
	}
	c.set(abstract, instance)
}

func (c *Container) remove(abstract string) {
	c.instances.Remove(abstract)
}

// Alias set alias for name
func (c *Container) Alias(abstract string, alias string) {
	if c.TryLock() {
		defer c.Unlock()
	}
	if abstract == alias {
		return
	}
	c.aliases.Set(alias, abstract)
}

func (c *Container) removeAlias(alias string) {
	c.aliases.Remove(alias)
}

// Extend extend
func (c *Container) Extend(abstract string, extenders ...func(instance any, container container.Container) any) {
	if c.TryLock() {
		defer c.Unlock()
	}
	abstract = c.realName(abstract)
	if instance, ok := c.instances.Get(abstract); ok {
		for _, extender := range extenders {
			instance = extender(instance, c)
		}
		c.instances.Set(abstract, instance)
	} else {
		if exts, ok := c.extends.Get(abstract); ok {
			exts.Push(extenders...)
		} else {
			c.extends.Set(abstract, lists.NewList(extenders...))
		}
	}
}
