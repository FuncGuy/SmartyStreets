package main

import "io"

type Pipeline struct {
	reader  io.ReadCloser
	writer  io.WriteCloser
	workers int

	verifier      Verifier
	verifyInput   chan *Envelope
	sequenceInput chan *Envelope
	writerInput   chan *Envelope
}

func Configure(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient, workers int) *Pipeline {
	return &Pipeline{
		reader:  reader,
		writer:  writer,
		workers: workers,

		verifier:      NewSmartyVerifier(client),
		verifyInput:   make(chan *Envelope, 1024),
		sequenceInput: make(chan *Envelope, 1024),
		writerInput:   make(chan *Envelope, 1024),
	}
}

func (this *Pipeline) Process() (err error) {

	this.startVerifyHandlers()

	// read from reader and write to verifyinput channel
	go func() {
		err = NewReaderHandler(this.reader, this.verifyInput).Handle()

	}()

	this.startSequenceHandler()
	this.awaitWriteHandler()

	return err

}

func (this *Pipeline) startVerifyHandlers() {
	for i := 0; i < this.workers; i++ {
		go NewVerifyHandler(this.verifyInput, this.sequenceInput, this.verifier).Handle()
	}
}

func (this *Pipeline) startSequenceHandler() {

	// read from reorder i/p and write to writeInput chanel
	go NewSequenceHandler(this.sequenceInput, this.writerInput).Handle() // runs in background
}

func (this *Pipeline) awaitWriteHandler() {
	// read from writer input channel and write to writer
	NewWriterHandler(this.writerInput, this.writer).Handle() // blocking
}
