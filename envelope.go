package main

const (
	initialSequenceValue = 0
	eofSequenceValue     = -1
)

var endOfFile = &Envelope{Sequence: eofSequenceValue}

type VerifyHandler struct {
	input    chan *Envelope
	output   chan *Envelope
	verifier Verifier
}

type Verifier interface {
	Verify(AddressInput) AddressOutput
}

type (
	Envelope struct {
		Input    AddressInput
		Output   AddressOutput
		Sequence int
	}

	AddressInput struct {
		Street1 string
		City    string
		State   string
		ZIPCode string
	}

	AddressOutput struct {
		Status        string
		DeliveryLine1 string
		LastLine      string
		City          string
		State         string
		ZIPCode       string
	}
)
