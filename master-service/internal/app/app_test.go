package app

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppRun(t *testing.T) {
	// t.Parallel() Not Parallelable test because of syscall.
	var (
		app *App
		err error
	)
	const (
		failTimeout = time.Second * 5

		// 1.5 sec.
		shutdownTimeout = time.Microsecond * 1500
	)
	errchan := make(chan error, 1)

	app, err = New()
	assert.NoError(t, err)

	shutdown := func() {
		time.Sleep(shutdownTimeout)
		err2 := syscall.Kill(os.Getpid(), syscall.SIGINT)
		assert.NoError(t, err2)
	}
	runner := func() {
		errchan <- app.Run()
		close(errchan)
	}

	go runner()
	go shutdown()

	timer := time.NewTimer(failTimeout)
	select {
	case <-timer.C:
		t.Fatalf("app shutdown has taken more than %s", failTimeout.String())
	case err = <-errchan:
		assert.NoError(t, err)
	}

}
