package dumbdumb

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime/debug"
)

// Listeners provide an interface for incoming requests to be submitted.  The
// listener sends requests to the incoming channel.
type RequestListener interface {
	Listen(incoming chan Request)
}

// Stores the payload of a request, and provides the mechanism for sending the
// response output back to the requestor.  Each request implementation is
// coupled with a particular Listener implementation.
type Request interface {
	GetPayload() (payload string)
	SendOutput(output string) (err error)
}

// Handlers take a request string, handle the request, and generate some output
// to be sent back to the requestor.
type Handler interface {
	HandleRequest(request Request) (err error)
}

// The Server is responsible for initializing a set of configured listeners
// which listen for incoming requests.  Requests are routed to configured
// handlers based on the payload of the request.
type Server struct {
	// TODO change to pointers?
	listeners []RequestListener
	handlers  map[string]Handler
}

func NewServer() *Server {
	return &Server{
		listeners: make([]RequestListener, 0),
		handlers:  make(map[string]Handler),
	}
}

func (s *Server) AddListener(l RequestListener) {
	s.listeners = append(s.listeners, l)
}

func (s *Server) AddHandler(pattern string, handler Handler) {
	s.handlers[pattern] = handler
}

func (s *Server) RouteRequest(request string) (*Handler, error) {
	for patt, h := range s.handlers {
		match, err := regexp.MatchString(patt, request)
		if err != nil {
			log.Printf("regexp match error: %v", err)
		}
		if match {
			return &h, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No handler matches request '%v'", request))
}

func (s *Server) HandleRequest(request Request) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			log.Printf("Unknown error caused panic! Error details: %v, Stacktrace: %s", r, stack)
		}
	}()
	req := request.GetPayload()
	handler, err := s.RouteRequest(req)
	if err != nil {
		log.Printf("Error routing request: '%v': %v", req, err)
		return
	}
	err = (*handler).HandleRequest(request)
	if err != nil {
		log.Printf("Error handling request '%v': %v", req, err)
		return
	}
}

// Begin listening for and handling incoming requests.
func (s *Server) ListenAndServe() error {
	ch := make(chan Request)
	for _, l := range s.listeners {
		go l.Listen(ch)
	}

	for {
		// handle the next request from incoming channel
		request := <-ch
		go s.HandleRequest(request)
	}
}
