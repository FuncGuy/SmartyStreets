package main

import (
	"encoding/csv"
	"errors"
	"io"
)

type ReaderHandler struct {
	reader   *csv.Reader
	closer   io.Closer
	output   chan *Envelope
	sequence int
	err      error
}

func NewReaderHandler(reader io.ReadCloser, output chan *Envelope) *ReaderHandler {
	return &ReaderHandler{
		reader:   csv.NewReader(reader),
		closer:   reader,
		output:   output,
		sequence: initialSequenceValue,
	}
}

func (this *ReaderHandler) Handle() error {
	defer this.close()
	// read from the reader
	// extract fields
	// create an envelope and send to output channel
	this.skipHeader() // skip the header before extracting actual contents
	//loop through all records until eof reached till then extract contents and create env and write to output.
	for {
		record, err := this.reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			this.err = err
			return errors.New("Malformed input")
		}
		this.sendEnvelope(record)

	}

	return nil
}

func (this *ReaderHandler) skipHeader() {
	this.reader.Read() // read till end of line (\n) read doc for more info.
}

func (this *ReaderHandler) sendEnvelope(record []string) {
	envelope := &Envelope{
		Sequence: this.sequence,
		Input:    createInput(record),
	}
	this.output <- envelope
	this.sequence++
}

func createInput(record []string) AddressInput {
	return AddressInput{
		Street1: record[0],
		City:    record[1],
		State:   record[2],
		ZIPCode: record[3],
	}
}

func (this *ReaderHandler) close() {
	if this.err == nil { // if error/empty is nil then only send to the output channel.
		this.output <- &Envelope{Sequence: this.sequence, EOF: true}
	}
	close(this.output)
	this.closer.Close()
}
