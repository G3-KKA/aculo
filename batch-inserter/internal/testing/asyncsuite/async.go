package asyncsuite

import (
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type asyncT struct {
	tt *testing.T
	mx sync.Mutex
}

func AsyncT(t *testing.T) *asyncT {
	return &asyncT{
		tt: t,
		mx: sync.Mutex{},
	}
}

func (t *asyncT) Log(args ...any) {
	t.mx.Lock()
	t.tt.Log(args...)
	t.mx.Unlock()
}

func (t *asyncT) Cleanup(f func()) {
	t.mx.Lock()
	t.tt.Cleanup(f)
	t.mx.Unlock()

}
func (t *asyncT) Deadline() (deadline time.Time, ok bool) {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.tt.Deadline()
}
func (t *asyncT) Error(args ...any) {
	t.mx.Lock()
	t.tt.Error(args...)
	t.mx.Unlock()

}
func (t *asyncT) Errorf(format string, args ...any) {
	t.mx.Lock()
	t.tt.Errorf(format, args...)
	t.mx.Unlock()

}
func (t *asyncT) Fail() {
	t.mx.Lock()
	t.tt.Fail()
	t.mx.Unlock()

}
func (t *asyncT) FailNow() {
	t.mx.Lock()
	t.tt.FailNow()
	t.mx.Unlock()
}
func (t *asyncT) Failed() bool {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.tt.Failed()

}
func (t *asyncT) Fatal(args ...any) {
	t.mx.Lock()
	t.tt.Fatal(args...)
	t.mx.Unlock()
}
func (t *asyncT) Fatalf(format string, args ...any) {
	t.mx.Lock()
	t.tt.Fatalf(format, args...)
	t.mx.Unlock()
}
func (t *asyncT) Helper() {
	t.mx.Lock()
	t.tt.Helper()
	t.mx.Unlock()

}
func (t *asyncT) Logf(format string, args ...any) {
	t.mx.Lock()
	t.tt.Logf(format, args...)
	t.mx.Unlock()
}
func (t *asyncT) Name() string {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.tt.Name()
}
func (t *asyncT) Parallel() {
	t.mx.Lock()
	t.tt.Parallel()
	t.mx.Unlock()
}
func (t *asyncT) Run(name string, f func(t *testing.T)) bool {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.tt.Run(name, f)
}
func (t *asyncT) Setenv(key string, value string) {
	t.mx.Lock()
	t.tt.Setenv(key, value)
	t.mx.Unlock()
}
func (t *asyncT) Skip(args ...any) {
	t.mx.Lock()
	t.tt.Skip(args...)
	t.mx.Unlock()
}
func (t *asyncT) SkipNow() {
	t.mx.Lock()
	t.tt.SkipNow()
	t.mx.Unlock()
}
func (t *asyncT) Skipf(format string, args ...any) {
	t.mx.Lock()
	t.tt.Skipf(format, args...)
	t.mx.Unlock()
}
func (t *asyncT) Skipped() bool {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.tt.Skipped()
}
func (t *asyncT) TempDir() string {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.tt.TempDir()
}

type asyncSuite struct {
	sut *suite.Suite
	t   *asyncT
	mx  sync.Mutex
}

func AsyncSuite(sut *suite.Suite) *asyncSuite {

	return &asyncSuite{
		sut: sut,
		mx:  sync.Mutex{},
		t:   AsyncT(sut.T()),
	}
}
func (asuite *asyncSuite) Require() *require.Assertions {
	panic("unimplemented func (asuite *asyncSuite) Require() *require.Assertions")
}

func (asuite *asyncSuite) Assert() *assert.Assertions {
	panic("unimplemented func (asuite *asyncSuite) Assert() *assert.Assertions ")
}

/*
asuite.mx.Lock()
asuite.sut.LOGIC
asuite.mx.Unlock()
*/

/*
asuite.mx.Lock()
defer asuite.mx.Unlock()
return asuite.sut.LOGIC
*/

func (asuite *asyncSuite) Condition(comp assert.Comparison, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Condition(comp, msgAndArgs...)
}
func (asuite *asyncSuite) Conditionf(comp assert.Comparison, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Conditionf(comp, msg, args...)
}
func (asuite *asyncSuite) Contains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Contains(s, contains, msgAndArgs...)
}
func (asuite *asyncSuite) Containsf(s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Containsf(s, contains, msg, args...)
}
func (asuite *asyncSuite) DirExists(path string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.DirExists(path, msgAndArgs...)
}
func (asuite *asyncSuite) DirExistsf(path string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.DirExistsf(path, msg, args...)
}
func (asuite *asyncSuite) ElementsMatch(listA interface{}, listB interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ElementsMatch(listA, listB, msgAndArgs...)
}
func (asuite *asyncSuite) ElementsMatchf(listA interface{}, listB interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ElementsMatchf(listA, listB, msg, args...)
}
func (asuite *asyncSuite) Empty(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Empty(object, msgAndArgs...)
}
func (asuite *asyncSuite) Emptyf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Emptyf(object, msg, args...)
}
func (asuite *asyncSuite) Equal(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Equal(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) EqualError(theError error, errString string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EqualError(theError, errString, msgAndArgs...)
}
func (asuite *asyncSuite) EqualErrorf(theError error, errString string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EqualErrorf(theError, errString, msg, args...)
}
func (asuite *asyncSuite) EqualExportedValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EqualExportedValues(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) EqualExportedValuesf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EqualExportedValuesf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) EqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EqualValues(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) EqualValuesf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EqualValuesf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) Equalf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Equalf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) Error(err error, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Error(err, msgAndArgs...)
}
func (asuite *asyncSuite) ErrorAs(err error, target interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ErrorAs(err, target, msgAndArgs...)
}
func (asuite *asyncSuite) ErrorAsf(err error, target interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ErrorAsf(err, target, msg, args...)
}
func (asuite *asyncSuite) ErrorContains(theError error, contains string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ErrorContains(theError, contains, msgAndArgs...)
}
func (asuite *asyncSuite) ErrorContainsf(theError error, contains string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ErrorContainsf(theError, contains, msg, args...)
}
func (asuite *asyncSuite) ErrorIs(err error, target error, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ErrorIs(err, target, msgAndArgs...)
}
func (asuite *asyncSuite) ErrorIsf(err error, target error, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.ErrorIsf(err, target, msg, args...)
}
func (asuite *asyncSuite) Errorf(err error, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Errorf(err, msg, args...)
}
func (asuite *asyncSuite) Eventually(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Eventually(condition, waitFor, tick, msgAndArgs...)
}
func (asuite *asyncSuite) EventuallyWithT(condition func(collect *assert.CollectT), waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EventuallyWithT(condition, waitFor, tick, msgAndArgs...)
}
func (asuite *asyncSuite) EventuallyWithTf(condition func(collect *assert.CollectT), waitFor time.Duration, tick time.Duration, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.EventuallyWithTf(condition, waitFor, tick, msg, args...)
}
func (asuite *asyncSuite) Eventuallyf(condition func() bool, waitFor time.Duration, tick time.Duration, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Eventuallyf(condition, waitFor, tick, msg, args...)
}
func (asuite *asyncSuite) Exactly(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Exactly(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) Exactlyf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Exactlyf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) Fail(failureMessage string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Fail(failureMessage, msgAndArgs...)
}
func (asuite *asyncSuite) FailNow(failureMessage string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.FailNow(failureMessage, msgAndArgs...)
}
func (asuite *asyncSuite) FailNowf(failureMessage string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.FailNowf(failureMessage, msg, args...)
}
func (asuite *asyncSuite) Failf(failureMessage string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Failf(failureMessage, msg, args...)
}
func (asuite *asyncSuite) False(value bool, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.False(value, msgAndArgs...)
}
func (asuite *asyncSuite) Falsef(value bool, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Falsef(value, msg, args...)
}
func (asuite *asyncSuite) FileExists(path string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.FileExists(path, msgAndArgs...)
}
func (asuite *asyncSuite) FileExistsf(path string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.FileExistsf(path, msg, args...)
}
func (asuite *asyncSuite) Greater(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Greater(e1, e2, msgAndArgs...)
}
func (asuite *asyncSuite) GreaterOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.GreaterOrEqual(e1, e2, msgAndArgs...)
}
func (asuite *asyncSuite) GreaterOrEqualf(e1 interface{}, e2 interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.GreaterOrEqualf(e1, e2, msg, args...)
}
func (asuite *asyncSuite) Greaterf(e1 interface{}, e2 interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Greaterf(e1, e2, msg, args...)
}
func (asuite *asyncSuite) HTTPBodyContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPBodyContains(handler, method, url, values, str, msgAndArgs...)
}
func (asuite *asyncSuite) HTTPBodyContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPBodyContainsf(handler, method, url, values, str, msg, args...)
}
func (asuite *asyncSuite) HTTPBodyNotContains(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPBodyNotContains(handler, method, url, values, str, msgAndArgs...)
}
func (asuite *asyncSuite) HTTPBodyNotContainsf(handler http.HandlerFunc, method string, url string, values url.Values, str interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPBodyNotContainsf(handler, method, url, values, str, msg, args...)
}
func (asuite *asyncSuite) HTTPError(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPError(handler, method, url, values, msgAndArgs...)
}
func (asuite *asyncSuite) HTTPErrorf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPErrorf(handler, method, url, values, msg, args...)
}
func (asuite *asyncSuite) HTTPRedirect(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPRedirect(handler, method, url, values, msgAndArgs...)
}
func (asuite *asyncSuite) HTTPRedirectf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPRedirectf(handler, method, url, values, msg, args...)
}
func (asuite *asyncSuite) HTTPStatusCode(handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPStatusCode(handler, method, url, values, statuscode, msgAndArgs...)
}
func (asuite *asyncSuite) HTTPStatusCodef(handler http.HandlerFunc, method string, url string, values url.Values, statuscode int, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPStatusCodef(handler, method, url, values, statuscode, msg, args...)
}
func (asuite *asyncSuite) HTTPSuccess(handler http.HandlerFunc, method string, url string, values url.Values, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPSuccess(handler, method, url, values, msgAndArgs...)
}
func (asuite *asyncSuite) HTTPSuccessf(handler http.HandlerFunc, method string, url string, values url.Values, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.HTTPSuccessf(handler, method, url, values, msg, args...)
}
func (asuite *asyncSuite) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Implements(interfaceObject, object, msgAndArgs...)
}
func (asuite *asyncSuite) Implementsf(interfaceObject interface{}, object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Implementsf(interfaceObject, object, msg, args...)
}
func (asuite *asyncSuite) InDelta(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InDelta(expected, actual, delta, msgAndArgs...)
}
func (asuite *asyncSuite) InDeltaMapValues(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InDeltaMapValues(expected, actual, delta, msgAndArgs...)
}
func (asuite *asyncSuite) InDeltaMapValuesf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InDeltaMapValuesf(expected, actual, delta, msg, args...)
}
func (asuite *asyncSuite) InDeltaSlice(expected interface{}, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InDeltaSlice(expected, actual, delta, msgAndArgs...)
}
func (asuite *asyncSuite) InDeltaSlicef(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InDeltaSlicef(expected, actual, delta, msg, args...)
}
func (asuite *asyncSuite) InDeltaf(expected interface{}, actual interface{}, delta float64, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InDeltaf(expected, actual, delta, msg, args...)
}
func (asuite *asyncSuite) InEpsilon(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InEpsilon(expected, actual, epsilon, msgAndArgs...)
}
func (asuite *asyncSuite) InEpsilonSlice(expected interface{}, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InEpsilonSlice(expected, actual, epsilon, msgAndArgs...)
}
func (asuite *asyncSuite) InEpsilonSlicef(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InEpsilonSlicef(expected, actual, epsilon, msg, args...)
}
func (asuite *asyncSuite) InEpsilonf(expected interface{}, actual interface{}, epsilon float64, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.InEpsilonf(expected, actual, epsilon, msg, args...)
}
func (asuite *asyncSuite) IsDecreasing(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsDecreasing(object, msgAndArgs...)
}
func (asuite *asyncSuite) IsDecreasingf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsDecreasingf(object, msg, args...)
}
func (asuite *asyncSuite) IsIncreasing(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsIncreasing(object, msgAndArgs...)
}
func (asuite *asyncSuite) IsIncreasingf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsIncreasingf(object, msg, args...)
}
func (asuite *asyncSuite) IsNonDecreasing(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsNonDecreasing(object, msgAndArgs...)
}
func (asuite *asyncSuite) IsNonDecreasingf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsNonDecreasingf(object, msg, args...)
}
func (asuite *asyncSuite) IsNonIncreasing(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsNonIncreasing(object, msgAndArgs...)
}
func (asuite *asyncSuite) IsNonIncreasingf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsNonIncreasingf(object, msg, args...)
}
func (asuite *asyncSuite) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsType(expectedType, object, msgAndArgs...)
}
func (asuite *asyncSuite) IsTypef(expectedType interface{}, object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.IsTypef(expectedType, object, msg, args...)
}
func (asuite *asyncSuite) JSONEq(expected string, actual string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.JSONEq(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) JSONEqf(expected string, actual string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.JSONEqf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) Len(object interface{}, length int, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Len(object, length, msgAndArgs...)
}
func (asuite *asyncSuite) Lenf(object interface{}, length int, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Lenf(object, length, msg, args...)
}
func (asuite *asyncSuite) Less(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Less(e1, e2, msgAndArgs...)
}
func (asuite *asyncSuite) LessOrEqual(e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.LessOrEqual(e1, e2, msgAndArgs...)
}
func (asuite *asyncSuite) LessOrEqualf(e1 interface{}, e2 interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.LessOrEqualf(e1, e2, msg, args...)
}
func (asuite *asyncSuite) Lessf(e1 interface{}, e2 interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Lessf(e1, e2, msg, args...)
}
func (asuite *asyncSuite) Negative(e interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Negative(e, msgAndArgs...)
}
func (asuite *asyncSuite) Negativef(e interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Negativef(e, msg, args...)
}
func (asuite *asyncSuite) Never(condition func() bool, waitFor time.Duration, tick time.Duration, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Never(condition, waitFor, tick, msgAndArgs...)
}
func (asuite *asyncSuite) Neverf(condition func() bool, waitFor time.Duration, tick time.Duration, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Neverf(condition, waitFor, tick, msg, args...)
}
func (asuite *asyncSuite) Nil(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Nil(object, msgAndArgs...)
}
func (asuite *asyncSuite) Nilf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Nilf(object, msg, args...)
}
func (asuite *asyncSuite) NoDirExists(path string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NoDirExists(path, msgAndArgs...)
}
func (asuite *asyncSuite) NoDirExistsf(path string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NoDirExistsf(path, msg, args...)
}
func (asuite *asyncSuite) NoError(err error, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NoError(err, msgAndArgs...)
}
func (asuite *asyncSuite) NoErrorf(err error, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NoErrorf(err, msg, args...)
}
func (asuite *asyncSuite) NoFileExists(path string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NoFileExists(path, msgAndArgs...)
}
func (asuite *asyncSuite) NoFileExistsf(path string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NoFileExistsf(path, msg, args...)
}
func (asuite *asyncSuite) NotContains(s interface{}, contains interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotContains(s, contains, msgAndArgs...)
}
func (asuite *asyncSuite) NotContainsf(s interface{}, contains interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotContainsf(s, contains, msg, args...)
}
func (asuite *asyncSuite) NotEmpty(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotEmpty(object, msgAndArgs...)
}
func (asuite *asyncSuite) NotEmptyf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotEmptyf(object, msg, args...)
}
func (asuite *asyncSuite) NotEqual(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotEqual(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) NotEqualValues(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotEqualValues(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) NotEqualValuesf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotEqualValuesf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) NotEqualf(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotEqualf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) NotErrorIs(err error, target error, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotErrorIs(err, target, msgAndArgs...)
}
func (asuite *asyncSuite) NotErrorIsf(err error, target error, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotErrorIsf(err, target, msg, args...)
}
func (asuite *asyncSuite) NotImplements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotImplements(interfaceObject, object, msgAndArgs...)
}
func (asuite *asyncSuite) NotImplementsf(interfaceObject interface{}, object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotImplementsf(interfaceObject, object, msg, args...)
}
func (asuite *asyncSuite) NotNil(object interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotNil(object, msgAndArgs...)
}
func (asuite *asyncSuite) NotNilf(object interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotNilf(object, msg, args...)
}
func (asuite *asyncSuite) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotPanics(f, msgAndArgs...)
}
func (asuite *asyncSuite) NotPanicsf(f assert.PanicTestFunc, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotPanicsf(f, msg, args...)
}
func (asuite *asyncSuite) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotRegexp(rx, str, msgAndArgs...)
}
func (asuite *asyncSuite) NotRegexpf(rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotRegexpf(rx, str, msg, args...)
}
func (asuite *asyncSuite) NotSame(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotSame(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) NotSamef(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotSamef(expected, actual, msg, args...)
}
func (asuite *asyncSuite) NotSubset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotSubset(list, subset, msgAndArgs...)
}
func (asuite *asyncSuite) NotSubsetf(list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotSubsetf(list, subset, msg, args...)
}
func (asuite *asyncSuite) NotZero(i interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotZero(i, msgAndArgs...)
}
func (asuite *asyncSuite) NotZerof(i interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.NotZerof(i, msg, args...)
}
func (asuite *asyncSuite) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Panics(f, msgAndArgs...)
}
func (asuite *asyncSuite) PanicsWithError(errString string, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.PanicsWithError(errString, f, msgAndArgs...)
}
func (asuite *asyncSuite) PanicsWithErrorf(errString string, f assert.PanicTestFunc, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.PanicsWithErrorf(errString, f, msg, args...)
}
func (asuite *asyncSuite) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.PanicsWithValue(expected, f, msgAndArgs...)
}
func (asuite *asyncSuite) PanicsWithValuef(expected interface{}, f assert.PanicTestFunc, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.PanicsWithValuef(expected, f, msg, args...)
}
func (asuite *asyncSuite) Panicsf(f assert.PanicTestFunc, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Panicsf(f, msg, args...)

}
func (asuite *asyncSuite) Positive(e interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Positive(e, msgAndArgs...)
}
func (asuite *asyncSuite) Positivef(e interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Positivef(e, msg, args...)
}
func (asuite *asyncSuite) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Regexp(rx, str, msgAndArgs...)
}
func (asuite *asyncSuite) Regexpf(rx interface{}, str interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Regexpf(rx, str, msg, args...)
}
func (asuite *asyncSuite) Run(name string, subtest func()) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Run(name, subtest)
}
func (asuite *asyncSuite) Same(expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Same(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) Samef(expected interface{}, actual interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Samef(expected, actual, msg, args...)
}
func (asuite *asyncSuite) SetS(s suite.TestingSuite) {
	asuite.mx.Lock()
	asuite.sut.SetS(s)
	asuite.mx.Unlock()
}
func (asuite *asyncSuite) SetT(t *testing.T) {
	asuite.mx.Lock()
	asuite.sut.SetT(t)
	asuite.t = AsyncT(t)
	asuite.mx.Unlock()
}
func (asuite *asyncSuite) Subset(list interface{}, subset interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Subset(list, subset, msgAndArgs...)
}
func (asuite *asyncSuite) Subsetf(list interface{}, subset interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Subsetf(list, subset, msg, args...)
}
func (asuite *asyncSuite) T() *asyncT {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()

	return AsyncT(asuite.sut.T())
}
func (asuite *asyncSuite) True(value bool, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.True(value, msgAndArgs...)
}
func (asuite *asyncSuite) Truef(value bool, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Truef(value, msg, args...)
}
func (asuite *asyncSuite) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.WithinDuration(expected, actual, delta, msgAndArgs...)
}
func (asuite *asyncSuite) WithinDurationf(expected time.Time, actual time.Time, delta time.Duration, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.WithinDurationf(expected, actual, delta, msg, args...)
}
func (asuite *asyncSuite) WithinRange(actual time.Time, start time.Time, end time.Time, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.WithinRange(actual, start, end, msgAndArgs...)
}
func (asuite *asyncSuite) WithinRangef(actual time.Time, start time.Time, end time.Time, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.WithinRangef(actual, start, end, msg, args...)
}
func (asuite *asyncSuite) YAMLEq(expected string, actual string, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.YAMLEq(expected, actual, msgAndArgs...)
}
func (asuite *asyncSuite) YAMLEqf(expected string, actual string, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.YAMLEqf(expected, actual, msg, args...)
}
func (asuite *asyncSuite) Zero(i interface{}, msgAndArgs ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Zero(i, msgAndArgs...)
}
func (asuite *asyncSuite) Zerof(i interface{}, msg string, args ...interface{}) bool {
	asuite.mx.Lock()
	defer asuite.mx.Unlock()
	return asuite.sut.Zerof(i, msg, args...)
}
