package logger

import (
	"os"
	"testing"
	"time"

	"master-service/config"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	const testDebugMessage = "Testing logger123."
	var (
		temp *os.File
		err  error

		logger Logger
	)

	temp, err = os.CreateTemp(os.TempDir(), "gotestfile*")
	assert.NoError(t, err)
	defer temp.Close()
	tempname := os.TempDir() + "/" + temp.Name()

	logger, err = New(config.Logger{
		SyncTimeout: time.Millisecond * 100,
		Cores: []config.LoggerCore{
			{
				Name:           "test",
				EncoderLevel:   "production",
				Path:           config.EnvString(tempname),
				Level:          -1,
				MustCreateCore: true,
			},
		},
	})
	assert.NoError(t, err)
	logger.Debug(testDebugMessage)
	bytes, err := os.ReadFile(tempname)
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), testDebugMessage)
}
