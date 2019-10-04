// verify_handler_test
package main

import (
	"strings"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestHandlerFixture(t *testing.T) {
	gunit.Run(new(HanlderFixture), t)
}

type HanlderFixture struct {
	*gunit.Fixture

	input       chan *Envelope
	output      chan *Envelope
	application *FakeVerifier
	handler     *VerifyHandler
}

func (this *HanlderFixture) Setup() {
	this.input = make(chan *Envelope, 10)
	this.output = make(chan *Envelope, 10)
	this.application = NewFakeVerifier()
	this.handler = NewVerifyHandler(this.input, this.output, this.application)
}

func (this *HanlderFixture) TestVerifierRecievesInput() {
	this.application.output = AddressOutput{DeliveryLine1: "DeliveryLine1"}

	envelope := this.enqueueEnvelope("street")

	close(this.input)

	// Handler recevies envelope from input channel and writes/assigns to application verifier and to output channel.
	this.handler.Handle()

	this.AssertEqual(envelope, <-this.output) // OK
	this.AssertEqual("STREET", envelope.Output.DeliveryLine1)
}

func (this *HanlderFixture) TestInputQueueDrained() {
	envelope1 := this.enqueueEnvelope("41")
	envelope2 := this.enqueueEnvelope("42")
	envelope3 := this.enqueueEnvelope("43")

	close(this.input)

	this.handler.Handle()

	this.AssertEqual(envelope1, <-this.output)
	this.AssertEqual(envelope2, <-this.output)
	this.AssertEqual(envelope3, <-this.output)

	this.AssertEqual("41", envelope1.Output.DeliveryLine1)

}

func (this *HanlderFixture) enqueueEnvelope(street1 string) *Envelope {

	envelope := &Envelope{Input: AddressInput{Street1: street1}} // create an envelope with AddressInput holds street1 value as 42

	this.input <- envelope // push envelope to input channel.

	return envelope
}

//////////////////////////////////////////////////////////

type FakeVerifier struct {
	input  AddressInput
	output AddressOutput
}

func NewFakeVerifier() *FakeVerifier {
	return &FakeVerifier{}
}

func (this *FakeVerifier) Verify(value AddressInput) AddressOutput {
	this.input = value
	return AddressOutput{DeliveryLine1: strings.ToUpper(value.Street1)}
}
