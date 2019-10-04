package main

type Verifier interface {
	Verify(AddressInput) AddressOutput
}
type VerifyHandler struct {
	input    chan *Envelope
	output   chan *Envelope
	verifier Verifier
}

func NewVerifyHandler(in, out chan *Envelope, verifier Verifier) *VerifyHandler {
	return &VerifyHandler{
		input:    in,
		output:   out,
		verifier: verifier,
	}
}

func (this *VerifyHandler) Handle() {

	for envelope := range this.input {

		envelope.Output = this.verifier.Verify(envelope.Input) // pass the input of channel to application verifier and get the modified output from the verifier.

		this.output <- envelope // and send the envelope to output channel.
	}
}
