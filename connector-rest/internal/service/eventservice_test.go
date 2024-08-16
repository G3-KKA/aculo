package service

import (
	"aculo/connector-restapi/internal/testutils"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(testSuite))
}
func (t *testSuite) SetupSuite() {
	testutils.DefaultSetup(t, "../../..")
}
