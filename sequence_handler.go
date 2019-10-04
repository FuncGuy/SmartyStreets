package main

type SequenceHandler struct {
	input   chan *Envelope
	output  chan *Envelope
	counter int
	buffer  map[int]*Envelope
}

func NewSequenceHandler(input, output chan *Envelope) *SequenceHandler {
	return &SequenceHandler{
		input:   input,
		output:  output,
		counter: initialSequenceValue,
		buffer:  make(map[int]*Envelope),
	}
}

func (this *SequenceHandler) Handle() {
	//this is intresting below code (this.output <- <- this.input) writes only one message to output
	//inorder to write all input messages to output it is necessary to use for loop
	//this.output <- <-this.input // recieve off of channel and send onto another channel.
	for envelope := range this.input {
		this.processEnvelope(envelope)
	}

	close(this.output)
}

func (this *SequenceHandler) processEnvelope(envelope *Envelope) {
	// Adding envelope sequence as key and its corresponding envelope example  [4,envelope -> 4], [2,envelope2]
	this.buffer[envelope.Sequence] = envelope
	this.sendBufferedEnvelopesInOrder()
}

func (this *SequenceHandler) sendBufferedEnvelopesInOrder() {
	for {
		envelope, found := this.buffer[this.counter]
		if !found {
			break
		}
		this.sendNextEnvelope(envelope)
	}
}

func (this *SequenceHandler) sendNextEnvelope(envelope *Envelope) {
	if envelope.EOF {
		close(this.input)
	} else {
		this.output <- envelope // send to the channel
	}
	delete(this.buffer, this.counter) // after sending to output channel clear it off so memory is enchanced
	this.counter++
}
