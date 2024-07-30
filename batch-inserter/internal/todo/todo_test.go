package todo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TodoTestSuite struct {
	suite.Suite
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
func (t *TodoTestSuite) Test_allocBatches() {
	allocBatches() // TODO
}
