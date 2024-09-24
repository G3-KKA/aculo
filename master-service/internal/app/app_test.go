package app

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppRun(t *testing.T) {
	// t.Parallel() Not Parallel test because of syscall
	var (
		app *App
		err error
	)

	interrupted := make(chan struct{})
	timer := time.NewTimer(time.Second * 4)

	go func() {
		select {
		case <-timer.C:
			t.Log("interrupt handler not working or app shutdown taken >4Sec")
			t.FailNow()
			os.Exit(1)
		case <-interrupted:
		}
	}()
	app, err = New()

	assert.NoError(t, err)
	go func() {
		time.Sleep(time.Second)
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()

	err = app.Run()

	assert.NoError(t, err)
	close(interrupted)
}
