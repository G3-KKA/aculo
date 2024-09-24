package logger

import (
	"master-service/internal/config"
	"os"
	"testing"
	"time"

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

	logger, err = New(config.Logger{
		SyncTimeout: time.Millisecond * 100,
		Cores: []config.LoggerCore{
			{
				Name:           "test",
				EncoderLevel:   "production",
				Path:           config.EnvString(os.TempDir() + "/" + temp.Name()),
				Level:          -1,
				MustCreateCore: true,
			},
		},
	})
	logger.Debug(testDebugMessage)
	bytes, err := os.ReadFile(os.TempDir() + "/" + temp.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), "Testing logger123.")
}
