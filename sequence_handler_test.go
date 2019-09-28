package main

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestSequenceHandler(t *testing.T) {
	gunit.Run(new(SequenceHandlerFixture), t)
}

type SequenceHandlerFixture struct {
	*gunit.Fixture

	input   chan *Envelope
	output  chan *Envelope
	handler *SequenceHandler
}

func (this *SequenceHandlerFixture) Setup() {
	this.input = make(chan *Envelope, 10)
	this.output = make(chan *Envelope, 10)
	this.handler = NewSequenceHandler(this.input, this.output)
}

func (this *SequenceHandlerFixture) TestExpectedEnvelopeSentToOutput() {
	this.sendEnvelopesInSequence(0, 1, 2, 3, 4) //send an envelope with valid sequence

	this.handler.Handle() // handle will read from input channel and write to output channel

	this.AssertSprintEqual([]int{0, 1, 2, 3, 4}, this.sequenceOrder()) // asserting envelope with the output envelope valid sequence
	this.AssertSprintEqual(map[int]Envelope{}, this.handler.buffer)    // make sure buffer is empty after the processing

}

func (this *SequenceHandlerFixture) TestEnvelopeRecievedOutOfOrder_BufferedUntilContiguousBlock() {
	this.sendEnvelopesInSequence(4, 2, 0, 3, 1) // unordered envelope seq

	this.handler.Handle()

	this.AssertSprintEqual([]int{0, 1, 2, 3, 4}, this.sequenceOrder()) // asserting for ordered seq
	this.AssertSprintEqual(map[int]Envelope{}, this.handler.buffer)    // make sure buffer is empty after the processing
}

func (this *SequenceHandlerFixture) sendEnvelopesInSequence(sequences ...int) {
	for _, sequence := range sequences {
		this.input <- &Envelope{Sequence: sequence}
	}
	close(this.input) // always after writing into the channel close the channel

}

//this mehtod will read sequence from the output channel and construct array of sequence
func (this *SequenceHandlerFixture) sequenceOrder() (order []int) {
	close(this.output) // always after reading from the channel close the channel
	for envelope := range this.output {
		order = append(order, envelope.Sequence)
	}
	return order

}
