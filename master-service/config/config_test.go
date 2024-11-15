package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
func (t *ConfigTestSuite) SetupSuite() {
	viper.Set("WORKSPACE", "..")
	viper.Set("CONFIG_FILE", viper.GetString("WORKSPACE")+"/config.yaml")
}

func (t *ConfigTestSuite) TestReadIn() {
	t.T().Parallel()
	cfg, err := Get()
	if t.NoError(err) {
		t.NotZero(cfg)
	}

}
func (t *ConfigTestSuite) TestHook() {
	t.T().Parallel()
	hook := envReplaceHook()
	t.NotNil(hook)
	// -1 represents any data that should not be parsed.
	testCases := []struct {
		// This kind of naming [F T D] used inside viper, this is not my fault.
		F              reflect.Type
		T              reflect.Type
		D              any
		ExpectedResult any
		Desc           string
	}{

		{
			F:              reflect.TypeOf(1),
			T:              reflect.TypeOf(nil),
			D:              -1,
			ExpectedResult: -1,
			Desc:           "hook input != reflect.String ",
		},
		{
			F:              reflect.TypeOf(EnvString("")),
			T:              reflect.TypeOf(nil),
			D:              -1,
			ExpectedResult: -1,
			Desc:           "hook target != reflect.config.config.Path ",
		},
		{
			F:              reflect.TypeOf(EnvString("")),
			T:              reflect.TypeOf(EnvString("")),
			D:              "/me/mario",
			ExpectedResult: "/me/mario",
			Desc:           "types are correct, but input missing ${ENV} expr",
		},
		{
			F:              reflect.TypeOf(EnvString("")),
			T:              reflect.TypeOf(EnvString("")),
			D:              "${WORKSPACE}/file/config.Path",
			ExpectedResult: viper.GetString("WORKSPACE") + "/file/config.Path",
			Desc:           "correct data, should be correct result",
		},
		{
			F:              reflect.TypeOf(EnvString("")),
			T:              reflect.TypeOf(EnvString("")),
			D:              "RR${WORKSPACE}/data",
			ExpectedResult: "RR" + viper.GetString("WORKSPACE") + "/data",
			Desc:           "correct data, should be correct result",
		},
		{
			F:              reflect.TypeOf(EnvString("")),
			T:              reflect.TypeOf(EnvString("")),
			D:              "RR${WORKSPACE}",
			ExpectedResult: "RR" + viper.GetString("WORKSPACE"),
			Desc:           "correct data, should be correct result",
		},
		{
			F:              reflect.TypeOf(EnvString("")),
			T:              reflect.TypeOf(EnvString("")),
			D:              "RR}${WORKSPACE",
			ExpectedResult: "RR}${WORKSPACE",
			Desc:           "",
		},
	}
	for _, testcase := range testCases {
		result, _ := hook(testcase.F, testcase.T, testcase.D)
		t.Equal(testcase.ExpectedResult, result, testcase.Desc)
	}
}
