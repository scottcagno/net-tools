package http

import (
	"net/http"
)

type Multiplexer struct {
	// dont forget handlers
	http.Handler
}

func NewMultiplexer() *Multiplexer {
	return &Multiplexer{}
}
