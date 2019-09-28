package main

import (
	"io"
)

type Handler interface {
}

type Wireup struct {
	reader io.ReadCloser
	writer io.WriteCloser
	client HTTPClient
}

func Configure(reader io.ReadCloser, writer io.WriteCloser, client HTTPClient) *Wireup {
	return &Wireup{
		reader: reader,
		writer: writer,
		client: client,
	}
}

func (this *Wireup) Build() Handler {
	return nil
}
