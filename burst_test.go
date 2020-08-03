package particles

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	burst3x30 = Burst{
		Count:    30,
		Cycles:   3,
		Interval: 0.10,
	}
)

func TestCalledTwice(t *testing.T) {
	b := burst3x30

	// Assert
	assert.Equal(t, 30, b.due())
	assert.Equal(t, 0, b.due())
}

func TestUpdateTriggersSecond(t *testing.T) {
	b := burst3x30
	assert.Equal(t, 30, b.due())
	assert.Equal(t, 0, b.due())

	// Act
	b.update(0.10)

	// Assert
	assert.Equal(t, 30, b.due())
}

func TestMultipleTimes(t *testing.T) {
	b := burst3x30

	// Act
	b.update(0.30)

	// Assert
	assert.Equal(t, 90, b.due())
	assert.Equal(t, 0, b.due())
}

func TestExhausted(t *testing.T) {
	b := burst3x30

	b.update(10.00)
	assert.Equal(t, 90, b.due())

	assert.Equal(t, 0, b.due())
}