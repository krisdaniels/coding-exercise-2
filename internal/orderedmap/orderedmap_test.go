package orderedmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	m := NewOrdererMap()
	m.Set("first", "some-value")
	m.Set("second", "some-value")
	m.Set("third", "some-value")

	assert.Equal(t, 3, len(m.(*orderedMap).kv))
}

func Test_Update(t *testing.T) {
	m := NewOrdererMap()
	m.Set("first", "some-value")
	m.Set("second", "some-value")
	m.Set("third", "some-value")

	m.Set("second", "some-other-value")

	res := m.Get("second")
	assert.Equal(t, "some-other-value", res)
}

func Test_DeleteFirst(t *testing.T) {
	m := NewOrdererMap()
	m.Set("first", "some-value")
	m.Set("second", "some-value")
	m.Set("third", "some-value")

	m.Delete("first")

	res := m.Get("first")
	assert.Equal(t, "", res)
	assert.Equal(t, "second", m.(*orderedMap).first.Key)
}

func Test_DeleteLast(t *testing.T) {
	m := NewOrdererMap()
	m.Set("first", "some-value")
	m.Set("second", "some-value")
	m.Set("third", "some-value")

	m.Delete("third")

	res := m.Get("third")
	assert.Equal(t, "", res)
	assert.Equal(t, "second", m.(*orderedMap).last.Key)
}

func Test_DeleteMiddle(t *testing.T) {
	m := NewOrdererMap()
	m.Set("first", "some-value")
	m.Set("second", "some-value")
	m.Set("third", "some-value")

	m.Delete("second")

	res := m.Get("second")
	assert.Equal(t, "", res)
	assert.Equal(t, "third", m.(*orderedMap).first.Next.Key)
	assert.Equal(t, "first", m.(*orderedMap).last.Prev.Key)
}

func Test_Iterate(t *testing.T) {
	m := NewOrdererMap()
	m.Set("key1", "val1")
	m.Set("key2", "val2")
	m.Set("key3", "val3")

	count := 0

	m.Iterate(func(key, value string) {
		count++
		assert.Equal(t, fmt.Sprintf("key%d", count), key)
		assert.Equal(t, fmt.Sprintf("val%d", count), value)
	})

	assert.Equal(t, 3, count)
}
