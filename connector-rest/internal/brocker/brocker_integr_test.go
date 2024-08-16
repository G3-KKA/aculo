package brocker

import (
	"aculo/connector-restapi/internal/testutils"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type BrockerIntegrTestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	viper.Set("INTEGRATION_TEST", true)
	suite.Run(t, new(BrockerIntegrTestSuite))
}

// ========================

func (t *BrockerIntegrTestSuite) SetupSuite() {
	testutils.DefaultSetup(t, "../..")
}
func (t *BrockerIntegrTestSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	}

}
func (t *BrockerIntegrTestSuite) Test_SendEvent() {

}
