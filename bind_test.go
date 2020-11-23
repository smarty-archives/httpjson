package httpjson

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBind_Successful(t *testing.T) {
	this := NewBindFixture(t)
	this.setRequestContentType("application/json")
	this.setRequestBody(`"hi"`)

	this.bind()

	this.assertSuccess()
	this.assertBoundModel("hi")
}
func TestBind_UnsupportedMediaType(t *testing.T) {
	this := NewBindFixture(t)
	this.setRequestContentType("text/plain")
	this.setRequestBody(`"hi"`)

	this.bind()

	this.assertFailure()
	this.assertBoundModel(nil)
	this.assertResultStatusCode(http.StatusUnsupportedMediaType)
	this.assertResultContentType("text/plain; charset=utf-8")
	this.assertResponseBodyPrefix("Unsupported Media Type (json content-type required)\n")
}
func TestBind_BadJSON(t *testing.T) {
	this := NewBindFixture(t)
	this.setRequestContentType("application/json")
	this.setRequestBody(`not valid json`)

	this.bind()

	this.assertFailure()
	this.assertBoundModel(nil)
	this.assertResultStatusCode(http.StatusBadRequest)
	this.assertResultContentType("text/plain; charset=utf-8")
	this.assertResponseBodyPrefix("Bad Request (json decode failure: ")
}

type BindFixture struct {
	t *testing.T

	request  *http.Request
	response *httptest.ResponseRecorder
	model    interface{}
	ok       bool
}

func NewBindFixture(t *testing.T) *BindFixture {
	t.Helper()
	t.Parallel()
	return &BindFixture{
		t:        t,
		request:  httptest.NewRequest("PUT", "/", nil),
		response: httptest.NewRecorder(),
	}
}
func (this *BindFixture) setRequestContentType(value string) {
	this.request.Header.Set("Content-Type", value)
}
func (this *BindFixture) setRequestBody(value string) {
	this.request.Body = ioutil.NopCloser(strings.NewReader(value))
}

func (this *BindFixture) bind() {
	this.ok = Bind(this.response, this.request, &this.model)
}

func (this *BindFixture) assertEqual(a, b interface{}) {
	this.t.Helper()
	assertEqual(this.t, a, b)
}
func (this *BindFixture) assertionReport(a, b interface{}) {
	this.t.Helper()
	assertionReport(this.t, a, b)
}
func (this *BindFixture) assertFailure() {
	this.t.Helper()
	this.assertEqual(false, this.ok)
}
func (this *BindFixture) assertSuccess() {
	this.t.Helper()
	this.assertEqual(true, this.ok)
}
func (this *BindFixture) assertBoundModel(expected interface{}) {
	this.t.Helper()
	this.assertEqual(expected, this.model)
}
func (this *BindFixture) assertResultStatusCode(expected int) {
	this.t.Helper()
	this.assertEqual(expected, this.response.Result().StatusCode)
}
func (this *BindFixture) assertResultContentType(expected string) {
	this.t.Helper()
	this.assertEqual(expected, this.response.Result().Header.Get("Content-Type"))
}
func (this *BindFixture) assertResponseBodyPrefix(expected string) {
	this.t.Helper()
	actual := this.responseBody()
	if !strings.HasPrefix(actual, expected) {
		this.assertionReport(expected, actual)
	}
}
func (this *BindFixture) responseBody() string {
	raw, _ := ioutil.ReadAll(this.response.Result().Body)
	return string(raw)
}
