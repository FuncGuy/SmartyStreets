package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/smartystreets/gunit"
)

func TestWriterHandlerFixture(t *testing.T) {
	gunit.Run(new(WriterHandlerFixture), t)
}

type WriterHandlerFixture struct {
	*gunit.Fixture

	//writer handler responsible for reading envelope from the input channel and writing to the csv.
	handler *WriterHandler
	//this is the input channel
	input chan *Envelope
	// buffer acts as a csv writer
	//in order to test we need to physcial file to avoid this we are creating spy buffer which accepts
	// some string and returns buffer string
	buffer *ReadWriteSpyBuffer
}

func (this *WriterHandlerFixture) Setup() {
	this.buffer = NewReadWriteSpyBuffer("") // Clear
	this.input = make(chan *Envelope, 10)   // Clear
	this.handler = NewWriterHandler(this.input, this.buffer)

}

var recordMatchingHeder = AddressOutput{
	Status:        "Status",
	DeliveryLine1: "DeliveryLine1",
	LastLine:      "LastLine",
	City:          "City",
	State:         "State",
	ZIPCode:       "ZIPCode",
}

func (this *WriterHandlerFixture) TestHeaderMatchesRecords() {
	this.input <- &Envelope{Output: recordMatchingHeder}

	close(this.input)

	this.handler.Handle()

	this.assertHeaderMatchesRecord()

}

func (this *WriterHandlerFixture) assertHeaderMatchesRecord() {
	lines := this.outputLines()

	header := lines[0]
	record := lines[1]

	this.AssertEqual(header, "Status,DeliveryLine1,LastLine,City,State,ZIPCode")
	this.AssertEqual(header, record)
}

//1) send envelope to the channel
//2) call handle function
//3) assert buffer content that contain envelope content
func (this *WriterHandlerFixture) TestEnvelopeWritten() {

	this.sendEnvelopes(1)

	this.handler.Handle()

	if lines := this.outputLines(); this.AssertEqual(2, len(lines)) {
		this.AssertEqual("A1,B1,C1,D1,E1,F1", lines[1])
	}

}

func (this *WriterHandlerFixture) TestAllEnvelopesWritten() {

	this.sendEnvelopes(2)

	this.handler.Handle()

	if lines := this.outputLines(); this.AssertEqual(3, len(lines)) {
		this.AssertEqual("A1,B1,C1,D1,E1,F1", lines[1])
		this.AssertEqual("A2,B2,C2,D2,E2,F2", lines[2])

	}
}

func (this *WriterHandlerFixture) TestOutputClosed() {
	close(this.input)
	this.handler.Handle()
	this.AssertEqual(1, this.buffer.closed)
}



func (this *WriterHandlerFixture) sendEnvelopes(count int) {
	for x := 1; x < count+1; x++ {
		this.input <- &Envelope{
			Output: createOutput(strconv.Itoa(x)),
		}
	}

	close(this.input)
}

func createOutput(index string) AddressOutput {

	return AddressOutput{
		Status:        "A" + index,
		DeliveryLine1: "B" + index,
		LastLine:      "C" + index,
		City:          "D" + index,
		State:         "E" + index,
		ZIPCode:       "F" + index,
	}

}

func (this *WriterHandlerFixture) outputLines() []string {

	outputFile := strings.TrimSpace(this.buffer.String())
	return strings.Split(outputFile, "\n")
}

/////////////////////////////////////////////////////////////////

type ReadWriteSpyBuffer struct {
	*bytes.Buffer
	closed int
}

func NewReadWriteSpyBuffer(value string) *ReadWriteSpyBuffer {
	return &ReadWriteSpyBuffer{
		Buffer: bytes.NewBufferString(value),
	}
}

func (this *ReadWriteSpyBuffer) Close() error {
	this.closed++
	return nil
}
