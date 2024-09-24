package config

import (
	"fmt"
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
	viper.Set("WORKSPACE", "../..")
	viper.Set("CONFIG_FILE", viper.GetString("WORKSPACE")+"/config.yaml")
}

func (t *ConfigTestSuite) TestReadIn() {
	cfg, err := ReadInConfig()
	if t.NoError(err) {
		t.NotZero(cfg)
	}

}
func (t *ConfigTestSuite) TestHook() {
	hook := envReplaceHook()
	t.NotNil(hook)
	// -1 represents any data that should not be parsed
	testCases := []struct {
		// This kind of naming [F T D] used inside viper, this is not my fault
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
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(nil),
			D:              -1,
			ExpectedResult: -1,
			Desc:           "hook target != reflect.config.config.Path ",
		},
		{
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(envstring("")),
			D:              "/me/mario",
			ExpectedResult: "/me/mario",
			Desc:           "hook input is config.config.Path type, but does not contain ${ENV} statement",
		},
		{
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(envstring("")),
			D:              "${WORKSPACE}/file/config.Path",
			ExpectedResult: viper.GetString("WORKSPACE") + "/file/config.Path",
			Desc:           "correct data, should be correct result",
		},
		{
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(envstring("")),
			D:              "RR${WORKSPACE}/data",
			ExpectedResult: "RR" + viper.GetString("WORKSPACE") + "/data",
			Desc:           "",
		},
		{
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(envstring("")),
			D:              "RR${WORKSPACE}",
			ExpectedResult: "RR" + viper.GetString("WORKSPACE"),
			Desc:           "",
		},
		{
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(envstring("")),
			D:              "${WORKSPACE}RR",
			ExpectedResult: viper.GetString("WORKSPACE") + "RR",
			Desc:           "",
		},
		{
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(envstring("")),
			D:              "RR}${WORKSPACE",
			ExpectedResult: "RR}${WORKSPACE",
			Desc:           "",
		},
		{
			F:              reflect.TypeOf(envstring("")),
			T:              reflect.TypeOf(envstring("")),
			D:              "RR}${",
			ExpectedResult: "RR}${",
			Desc:           "",
		},
	}
	for _, testcase := range testCases {
		result, _ := hook(testcase.F, testcase.T, testcase.D)
		t.Equal(testcase.ExpectedResult, result, testcase.Desc)
	}
}
func (t *ConfigTestSuite) Test_extFromPath() {
	testCases := []struct {
		Path  string
		Exted string
	}{
		{
			Path:  "some/config.yaml",
			Exted: "yaml",
		},

		{
			Path:  "config.json",
			Exted: "json",
		},
	}
	for _, testcase := range testCases {
		ext := extFromPath(testcase.Path)
		t.Equal(testcase.Exted, ext)
	}
}

func (t *ConfigTestSuite) Test_registerENV() {
	testCases := []struct {
		ENV    string
		Result string
		Error  error

		Desc string
	}{
		{
			ENV:    "WORKSPACE",
			Result: viper.GetString("WORKSPACE"),
			Error:  nil,
			Desc:   "WORKSPACE should be correct",
		},
		{
			ENV:    "Undefined",
			Result: "",
			Error:  fmt.Errorf("some error"),
			Desc:   "Undefined should be empty",
		},
	}

	for _, testcase := range testCases {
		err := registerENV(testcase.ENV)
		if t.Equal(testcase.Result, viper.GetString(testcase.ENV), testcase.Desc) {
			continue
		}
		t.ErrorIs(err, testcase.Error, testcase.Desc)
	}

}

func (t *ConfigTestSuite) Test_InitConfig() {
	err := initConfig()
	t.NoError(err, "should be ok")
}
