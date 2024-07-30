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

// ========================

func (t *ConfigTestSuite) SetupSuite() {
	viper.Set("WORKSPACE", "../..")
	viper.Set("CONFIG_FILE", viper.GetString("WORKSPACE")+"/config.yaml")
}
func (t *ConfigTestSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "Test_fillGlobalConfig":
		handleConfigFile()
	case "Test_registerENV":
	}

}
func (t *ConfigTestSuite) Test_execute() {
	defer func() {
		err := recover()
		t.Nil(err, "should not panic", err)
	}()
	execute([]initPhase{})
}
func (t *ConfigTestSuite) Test_bindFlags() {
	err := bindFlags()
	t.Nil(err, "should be able to bind flags", err)

}
func (t *ConfigTestSuite) Test_fillGlobalConfig() {
	err := fillGlobalConfig()
	t.Nil(err, "should be able to fill global config", err)

}
func (t *ConfigTestSuite) Test_envReplaceHook() {
	hook := envReplaceHook()
	t.NotNil(hook, "should be able to get hook,-")
	// -1 represents any data that should not be parsed
	testCases := []struct {
		// This kind of naming [F T D] used inside viper, this is not my fault
		F              reflect.Type
		T              reflect.Type
		D              any
		ExpectedResult any
		ExpecteedErr   error
		Desc           string
	}{

		{
			F:              reflect.TypeOf(1),
			T:              reflect.TypeOf(nil),
			D:              -1,
			ExpectedResult: -1,
			ExpecteedErr:   nil,
			Desc:           "hook input != reflect.String ",
		},
		{
			F:              reflect.TypeOf(path("")),
			T:              reflect.TypeOf(nil),
			D:              -1,
			ExpectedResult: -1,
			ExpecteedErr:   nil,
			Desc:           "hook target != reflect.config.path ",
		},
		{
			F:              reflect.TypeOf(path("")),
			T:              reflect.TypeOf(path("")),
			D:              "/me/mario",
			ExpectedResult: "/me/mario",
			ExpecteedErr:   nil,
			Desc:           "hook input is config.path type, but does not contain ${ENV} statement",
		},
		{
			F:              reflect.TypeOf(path("")),
			T:              reflect.TypeOf(path("")),
			D:              "${WORKSPACE}/file/path",
			ExpectedResult: viper.GetString("WORKSPACE") + "/file/path",
			ExpecteedErr:   nil,
			Desc:           "correct data, should be correct result",
		},
	}
	for _, testcase := range testCases {
		result, err := hook(testcase.F, testcase.T, testcase.D)
		t.Equal(testcase.ExpectedResult, result, testcase.Desc)
		t.Equal(testcase.ExpecteedErr, err, testcase.Desc)
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
		t.NotNil(testcase.Error, err, testcase.Desc)
	}

}

func (t *ConfigTestSuite) Test_InitConfig() {
	err := InitConfig()
	t.Equal(nil, err, "should be ok")
}
