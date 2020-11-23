package httpjson

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestErrorImplementErrorInterfaceButIsNotImplemented(t *testing.T) {
	var recovered interface{}
	defer func() {
		recovered = recover()
		if recovered != notImplemented {
			assertionReport(t, notImplemented, recovered)
		}
	}()

	var err error = Error{}
	_ = err.Error()
}
func TestRenderWithZeroOptions(t *testing.T) {
	this := NewRenderFixture(t)
	this.render()

	this.assertStatusCode(http.StatusOK)
	this.assertHeader("Content-Type", "application/json; charset=utf-8")
	this.assertBody("{}\n")
}
func TestRenderStatusCode(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderStatusCode(http.StatusTeapot),
	)
	this.assertStatusCode(http.StatusTeapot)
}
func TestRenderHeader(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderHeader("key", "value1", "value2"),
	)
	this.assertHeader("Key", "value1", "value2")
}
func TestRenderHeader_ContentType_Nop(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderHeader("Content-Type", "nope"),
	)
	this.assertHeader("Content-Type", "application/json; charset=utf-8")
}
func TestRenderData(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderIndent("", ""),
		Options.RenderData("this is a test"),
	)
	this.assertBody(`{"data":"this is a test"}` + "\n")
}
func TestRenderDataIndent(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderIndent(">>", "  "),
		Options.RenderData("this is a test"),
	)
	this.assertBody("{\n>>  \"data\": \"this is a test\"\n>>}\n")
}
func TestRenderInvalidData(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderData(make(chan int)),
	)
	this.assertStatusCode(http.StatusInternalServerError)
}
func TestErrorSerialization(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderIndent("", ""),
		Options.RenderMappedErrors(this.MapError, errors.New("1"), errors.New("2")),
	)
	this.assertBody(`{"errors":[{"id":"i:1","message":"m:1","fields":["f:1","f:1"]},{"id":"i:2","message":"m:2","fields":["f:2","f:2"]}]}` + "\n")
}
func TestFirstErrorStatusCodeUsedAsStatusCode(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderMappedErrors(this.MapError, errors.New("1"), errors.New("2")),
	)
	this.assertStatusCode(1)
}
func TestStatusCodeOverridesErrors(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderStatusCode(http.StatusTeapot),
		Options.RenderMappedErrors(this.MapError, errors.New("1"), errors.New("2")),
	)
	this.assertStatusCode(http.StatusTeapot)
}
func TestInternalServerErrorErrorOverridesEverything(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderStatusCode(http.StatusTeapot),
		Options.RenderErrors(Error{StatusCodeValue: http.StatusOK}),
		Options.RenderErrors(Error{StatusCodeValue: http.StatusInternalServerError}),
	)
	this.assertStatusCode(http.StatusInternalServerError)
}
func TestAboveInternalServerErrorErrorOverridesEverything(t *testing.T) {
	this := NewRenderFixture(t)
	this.render(
		Options.RenderStatusCode(http.StatusTeapot),
		Options.RenderErrors(Error{StatusCodeValue: http.StatusOK}),
		Options.RenderErrors(Error{StatusCodeValue: http.StatusBadGateway}),
	)
	this.assertStatusCode(http.StatusBadGateway)
}

func (this *RenderFixture) MapError(err error) Error {
	text := err.Error()
	parsed, _ := strconv.Atoi(text)
	return Error{
		StatusCodeValue: parsed,
		ID:              "i:" + text,
		Message:         "m:" + text,
		Fields:          []string{"f:" + text, "f:" + text},
	}
}

type RenderFixture struct {
	t *testing.T

	recorder *httptest.ResponseRecorder
	response *http.Response
	body     string
}

func NewRenderFixture(t *testing.T) *RenderFixture {
	return &RenderFixture{
		t: t,

		recorder: httptest.NewRecorder(),
	}
}

func (this *RenderFixture) render(options ...RenderOption) {
	Render(this.recorder, options...)
	this.response = this.recorder.Result()
	this.body = this.readBody()
}
func (this *RenderFixture) readBody() string {
	all, err := ioutil.ReadAll(this.response.Body)
	if err != nil {
		this.t.Fatal(err)
	}
	return string(all)
}

func (this *RenderFixture) assertEqual(a, b interface{}) {
	this.t.Helper()
	assertEqual(this.t, a, b)
}
func (this *RenderFixture) assertionReport(a, b interface{}) {
	this.t.Helper()
	assertionReport(this.t, a, b)
}
func (this *RenderFixture) assertStatusCode(expected int) {
	this.t.Helper()
	this.assertEqual(expected, this.response.StatusCode)
}
func (this *RenderFixture) assertHeader(key string, values ...string) {
	this.t.Helper()
	this.assertEqual(values, this.response.Header[key])
}
func (this *RenderFixture) assertBody(expected string) {
	this.t.Helper()
	this.assertEqual(expected, this.body)
}

func assertEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		assertionReport(t, a, b)
	}
}
func assertionReport(t *testing.T, a, b interface{}) {
	t.Helper()
	t.Errorf("\nExpected: %v\nActual:   %v", a, b)
}
