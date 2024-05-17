package container

import (
	"reflect"
	"testing"

	"github.com/gopi-frame/contract/container"
	"github.com/stretchr/testify/assert"
)

type _server001 struct {
	id   string
	name string
}

func (s _server001) Build(container container.Container) any {
	return &_server001{id: "_server_001", name: "testServer"}
}

type _server001Extend struct{}

func (s _server001Extend) Extend(instance any, c container.Container) any {
	instance.(*_server001).id = "0010"
	return instance
}

type _server002 struct {
	id   string
	name string
}

func (s _server002) Build(container container.Container) any {
	return &_server002{}
}

func (s _server002) Extend(instance any, container container.Container) any {
	instance.(*_server002).id = "_server002"
	return instance
}

func TestContainer(t *testing.T) {
	c := NewContainer()
	c.Set("basepath", "/path/to/app")
	assert.True(t, c.Has("basepath"))
	assert.Equal(t, "/path/to/app", c.Get("basepath"))

	c.Bind("server", _server001{}.Build)
	c.Alias("server", "server_001")
	assert.True(t, c.Has("server"))
	c.BindIf("server", _server002{}.Build)
	assert.Equal(t, reflect.TypeOf(new(_server001)).Elem(), reflect.TypeOf(c.Get("server")).Elem())
	assert.Equal(t, c.Get("server"), c.Get("server_001"))

	c.Extend("server_001", _server001Extend{}.Extend)
	assert.Equal(t, "0010", c.Get("server").(*_server001).id)
}
