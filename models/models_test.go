package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_ToJSONString(t *testing.T) {
	msg := Message{
		ID:        "8888-8888-8888-8888",
		Value:     "hello world",
		Timestamp: 1555566556,
	}
	expectedStr := `{"id":"8888-8888-8888-8888","value":"hello world","timestamp":1555566556}`
	res, err := msg.ToJSONString()
	assert.NoError(t, err)
	assert.Equal(t, expectedStr, res)
}
