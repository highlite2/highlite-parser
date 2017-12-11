package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemo_GetKey(t *testing.T) {
	// arrange
	m := NewMemo()
	// act
	res, err := m.Get("test", func() (interface{}, error) {
		return "ok", nil
	})
	val, _ := res.(string)
	// assert
	assert.Nil(t, err)
	assert.Equal(t, "ok", val)
}

func TestMemo_GetKeySeveralTimes(t *testing.T) {
	// arrange
	count := 0
	callback := func() (interface{}, error) {
		count++
		return "ok", nil
	}
	m := NewMemo()
	// act
	m.Get("test1", callback)
	m.Get("test2", callback)
	m.Get("test2", callback)
	m.Get("test2", callback)
	m.Get("test2", callback)
	res, err := m.Get("test2", callback)
	val, _ := res.(string)
	// assert
	assert.Equal(t, 2, count)
	assert.Nil(t, err)
	assert.Equal(t, "ok", val)
}

func TestMemo_GetKeyCallbackError(t *testing.T) {
	// arrange
	count := 0
	callback := func() (interface{}, error) {
		count++
		if count == 1 {
			return nil, fmt.Errorf("error")
		}

		return "second call", nil
	}
	m := NewMemo()
	// act
	r1, e1 := m.Get("test", callback)
	r2, e2 := m.Get("test", callback)
	v2, _ := r2.(string)
	// arrange
	assert.Nil(t, r1)
	assert.EqualError(t, e1, "error")
	assert.Nil(t, e2)
	assert.Equal(t, "second call", v2)
	assert.Equal(t, 2, count)
}
