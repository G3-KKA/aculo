package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppRun(t *testing.T) {

	var (
		app *App
		err error
	)

	app, err = New()

	assert.NoError(t, err)

	err = app.Run()

	assert.NoError(t, err)
}
