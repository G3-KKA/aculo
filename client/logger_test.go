package aculo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHappyPath(t *testing.T) {
	t.Parallel()

	const (
		message = `{"msg": "test succesful"}`
		address = "localhost:7730"
	)

	ctx := context.TODO()
	logger, err := New(ctx, address)
	assert.NoError(t, err)

	written, err := logger.Write([]byte(message))
	assert.NoError(t, err)
	assert.Len(t, message, written)

	err = logger.Close()
	assert.NoError(t, err)
}
