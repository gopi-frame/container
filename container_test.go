package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer(t *testing.T) {
	c := New[int]()
	c.Set("value1", 1)
	assert.True(t, c.Has("value1"))
	assert.Equal(t, 1, c.Get("value1"))

	var value = 2
	c.Defer("value2", func() (int, error) {
		defer func() {
			value++
		}()
		return value, nil
	})
	assert.True(t, c.Has("value2"))
	assert.Equal(t, 2, c.Get("value2"))
	assert.Equal(t, 3, value)
	assert.Equal(t, 2, c.Get("value2"))
	v, err := c.Make("value2")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, 3, v)
	assert.Equal(t, 4, value)
	_, err = c.Make("value3")
	assert.Error(t, err)

	v = c.Get("value3")
	assert.Equal(t, 0, v)
	v, err = c.GetE("value3")
	assert.Error(t, err)
	assert.Equal(t, 0, v)
}
