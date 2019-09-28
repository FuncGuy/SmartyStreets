package main

import (
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestWireUpFixture(t *testing.T) {
	gunit.Run(new(WireUpFixture), t)
}

type WireUpFixture struct {
	*gunit.Fixture
	reader  *ReadWriteSpyBuffer
	writer  *ReadWriteSpyBuffer
	client  *FakeHttpClient
	handler Handler
}

func (this *WireUpFixture) Setup() {
	this.reader = NewReadWriteSpyBuffer("")
	this.writer = NewReadWriteSpyBuffer("")
	this.client = &FakeHttpClient{}
	this.handler = Configure(this.reader, this.writer, this.client).Build()

}

func (this *WireUpFixture) Test() {
	this.client.Configure(rawJSONOutput, http.StatusOK, nil)
	this.reader.WriteString("Street1,City,State,ZIPCode")
	this.reader.WriteString("A,B,C,D")

	this.handler.Handle()

	this.So(this.writer.String(), should.Equal, "")

}
