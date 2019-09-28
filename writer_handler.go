package main

import (
	"encoding/csv"
	"io"
)

type WriterHandler struct {
	input  chan *Envelope
	closer io.Closer
	writer *csv.Writer
}

// function accepts input channel [chan * envelope] and output [io.writeCloser]
// writes the content of envelope to writer
// second parameter is Writecloser instead of writer because once the write is done
// writer has to close cleanly.
func NewWriterHandler(input chan *Envelope, output io.WriteCloser) *WriterHandler {
	this := &WriterHandler{
		input:  input,
		closer: output,                // io.closer and io.WriteCloser are from same struct so it is OK...
		writer: csv.NewWriter(output), // if doubt see the inner impl of csv.NewWriter
	}
	this.writeValues("Status", "DeliveryLine1", "LastLine", "City", "State", "ZIPCode")
	return this
}

func (this *WriterHandler) Handle() {

	//loop until the channel is closed
	for envelope := range this.input {
		this.writeAddressOutput(envelope.Output)
	}
	this.writer.Flush()
	this.closer.Close()

}

func (this *WriterHandler) writeAddressOutput(output AddressOutput) {
	this.writeValues(
		output.Status,
		output.DeliveryLine1,
		output.LastLine,
		output.City,
		output.State,
		output.ZIPCode,
	)
}

func (this *WriterHandler) writeValues(values ...string) {
	this.writer.Write(values)
}
