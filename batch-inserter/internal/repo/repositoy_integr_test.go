package repository

import (
	"aculo/batch-inserter/internal/testutils"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepoTestSuite_Integr struct {
	suite.Suite
}

func TestRepo_Integr(t *testing.T) {
	suite.Run(t, new(RepoTestSuite_Integr))
}

// ========================

func (t *RepoTestSuite_Integr) SetupSuite() {
	testutils.ThisIsIntegrationTest(t)
}
func (t *RepoTestSuite_Integr) BeforeTest(suiteName, testName string) {
	switch testName {
	}

}
