package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestAuthenticationClient(t *testing.T) {
	gunit.Run(new(AuthenticationClientFixture), t)
}

type AuthenticationClientFixture struct {
	*gunit.Fixture

	inner  *FakeHttpClient
	client *AuthenticationClient
}

func (this *AuthenticationClientFixture) Setup() {
	this.inner = &FakeHttpClient{}
	this.client = NewAuthenticationClient(this.inner, "https", "different-company.com", "authid", "authtoken")

}

// 1) verify app will send the "Request" with all query parameters.
// 2) Here a new client should do the authentication tasks of the main request.
func (this *AuthenticationClientFixture) TestHostnameAndSchemeAddedBeforeRequestIsSent() {
	request := httptest.NewRequest("GET", "/path?existingKey=existingValue", nil)

	this.client.Do(request)

	this.assertRequestConnectionInformation()
	this.assertQueryStringIncludesAuthentication()
}
func (this *AuthenticationClientFixture) assertRequestConnectionInformation() {
	this.AssertEqual("https", this.inner.request.URL.Scheme)
	this.AssertEqual("different-company.com", this.inner.request.Host)
	this.AssertEqual("different-company.com", this.inner.request.URL.Host)
}
func (this *AuthenticationClientFixture) assertQueryStringIncludesAuthentication() {
	this.assertQueryStringValue("auth-id", "authid")
	this.assertQueryStringValue("auth-token", "authtoken")
	this.assertQueryStringValue("existingKey", "existingValue")
}
func (this *AuthenticationClientFixture) assertQueryStringValue(key, expectedValue string) {

	this.AssertEqual(expectedValue, this.inner.request.URL.Query().Get(key))

}
func (this *AuthenticationClientFixture) TestResponseAndErrorFromInnerClientReturned() {

	this.inner.response = &http.Response{StatusCode: http.StatusTeapot}

	this.inner.err = errors.New("HTTP Error")

	request := httptest.NewRequest("GET", "/path", nil)

	response, err := this.client.Do(request)

	this.AssertEqual(response.StatusCode, http.StatusTeapot)
	this.AssertEqual(err.Error(), "HTTP Error")
}
