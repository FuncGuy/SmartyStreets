package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestVerifierFixture(t *testing.T) {
	gunit.Run(new(VerifierFixture), t)
}

type VerifierFixture struct {
	*gunit.Fixture

	client   *FakeHttpClient // client to perfom http request
	verifier *SmartyVerifier // verifier which takes input address and sends to smartyVerifier
}

func (this *VerifierFixture) Setup() {
	this.client = &FakeHttpClient{}                // stub the client
	this.verifier = NewSmartyVerifier(this.client) // stub the verifier by passing dependent http client
}

func (this *VerifierFixture) TestRequestComposedProperly() {
	input := AddressInput{
		Street1: "Street1",
		City:    "City",
		State:   "State",
		ZIPCode: "ZIPCode",
	}

	this.client.Configure("[{}]", http.StatusOK, nil)

	this.verifier.Verify(input)

	this.AssertEqual("GET", this.client.request.Method)
	this.AssertEqual("/street-address", this.client.request.URL.Path)
	this.AssertQueryStringValue("street", input.Street1)
	this.AssertQueryStringValue("city", input.City)
	this.AssertQueryStringValue("state", input.State)
	this.AssertQueryStringValue("zipCode", input.ZIPCode)

}

func (this *VerifierFixture) TestResponseParseRead() {
	this.client.Configure(rawJSONOutput, http.StatusOK, nil)
	result := this.verifier.Verify(AddressInput{})
	this.AssertEqual("1 Santa Claus Ln", result.DeliveryLine1)
	this.AssertEqual("North Pole AK 99705-9901", result.LastLine)
	this.AssertEqual("North Pole", result.City)
	this.AssertEqual("AK", result.State)
	this.AssertEqual("99705", result.ZIPCode)

}

func (this *VerifierFixture) TestMalformedJSONHandled() {
	const malformedJSONOutput = `alert('Hello World!', DROP TABLE Users);`
	this.client.Configure(malformedJSONOutput, http.StatusOK, nil)
	result := this.verifier.Verify(AddressInput{})
	this.AssertEqual("Invalid API Response", result.Status)
}

func (this *VerifierFixture) TestHTTPErrorHandled() {
	this.client.Configure("", 0, errors.New("GOPHERS!"))
	result := this.verifier.Verify(AddressInput{})
	this.AssertEqual("Invalid API Response", result.Status)
}

func (this *VerifierFixture) TestHTTPResponseBodyClosed() {
	this.client.Configure(rawJSONOutput, http.StatusOK, nil)
	this.verifier.Verify(AddressInput{})
	this.AssertEqual(1, this.client.responseBody.closed)
}

func (this *VerifierFixture) TestAddressStatus() {
	var (
		deliverableJSON      = buildAnalysisJSON("Y", "N", "Y")
		missingSecondaryJSON = buildAnalysisJSON("D", "N", "Y")
		droppedSecondaryJSON = buildAnalysisJSON("S", "N", "Y")
		vacantJSON           = buildAnalysisJSON("Y", "Y", "Y")
		inactiveJSON         = buildAnalysisJSON("Y", "N", "N")
		invalidJSON          = buildAnalysisJSON("N", "?", "?")
	)

	this.verifyAndAssertStatus(deliverableJSON, "Deliverable")
	this.verifyAndAssertStatus(missingSecondaryJSON, "Deliverable")
	this.verifyAndAssertStatus(droppedSecondaryJSON, "Deliverable")
	this.verifyAndAssertStatus(vacantJSON, "Vacant")
	this.verifyAndAssertStatus(inactiveJSON, "Inactive")
	this.verifyAndAssertStatus(invalidJSON, "Invalid")

}

func (this *VerifierFixture) verifyAndAssertStatus(jsonResponse, expectedStatus string) {

	this.client.Configure(jsonResponse, http.StatusOK, nil)
	output := this.verifier.Verify(AddressInput{})
	this.AssertEqual(expectedStatus, output.Status)

}

func buildAnalysisJSON(match, vacant, active string) string {
	template := `
		[
			{
		        "analysis": {
		            "dpv_match_code":"%s",
		            "dpv_vacant":"%s",
		            "active":"%s"        
			 	}
		    }
		]`

	return fmt.Sprintf(template, match, vacant, active)
}

const rawJSONOutput = `
[ 
	{  
		
	    "delivery_line_1": "1 Santa Claus Ln", 
	    "last_line": "North Pole AK 99705-9901",
		"components": { 
			"city_name": "North Pole",
			"state_abbreviation": "AK",
			"zipcode": "99705"
		}
	}
]`

func (this *VerifierFixture) AssertQueryStringValue(key, expected string) {
	query := this.client.request.URL.Query()
	this.AssertEqual(expected, query.Get(key))
}

func (this *VerifierFixture) rawQuery() string {
	return this.client.request.URL.RawQuery
}

//////////////////////////////////////////////

type FakeHttpClient struct {
	request      *http.Request
	response     *http.Response
	responseBody *VerifierSpyBuffer
	err          error
}

func NewFakeHTTPClient(responseText string, StatusCode int, err error) *FakeHttpClient {
	return &FakeHttpClient{
		response: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(responseText)),
			StatusCode: StatusCode,
		},
		err: err,
	}
}

func (this *FakeHttpClient) Configure(responseText string, statusCode int, err error) {

	if err == nil {
		this.responseBody = NewVerifierSpyBuffer(responseText)
		this.response = &http.Response{
			Body:       this.responseBody,
			StatusCode: statusCode,
		}
	}
	this.err = err

}

func (this *FakeHttpClient) Do(request *http.Request) (*http.Response, error) {
	this.request = request
	return this.response, this.err
}

//////////////////////////////////////////////

type VerifierSpyBuffer struct {
	*bytes.Buffer
	closed int
}

func NewVerifierSpyBuffer(value string) *VerifierSpyBuffer {
	return &VerifierSpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *VerifierSpyBuffer) Close() error {
	this.closed++
	this.Buffer.Reset()
	return nil
}
