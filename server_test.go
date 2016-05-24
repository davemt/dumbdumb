package dumbdumb

import (
	"testing"
)

type MockRequest struct {
	Payload string
}

func (r MockRequest) GetPayload() string { return r.Payload }

func (r MockRequest) SendOutput(output string) error { return nil }

type MockHandler struct {
	Name string
}

func (h MockHandler) HandleRequest(request Request) error { return nil }

func TestAddHandler(t *testing.T) {
	server := NewServer()

	server.AddHandler("requesttype1.*", MockHandler{})
	server.AddHandler("requesttype2.*", MockHandler{})

	if len(server.handlers) != 2 {
		t.Error("Expected to have 2 handlers added to server, had",
			len(server.handlers))
	}
}

func TestRouteRequest(t *testing.T) {
	server := NewServer()

	h1 := MockHandler{Name: "handler1"}
	h2 := MockHandler{Name: "handler2"}
	server.AddHandler("requesttype1.*", h1)
	server.AddHandler("requesttype2.*", h2)

	h, err := server.RouteRequest("requesttype1 blah blah")
	if err != nil {
		t.Fatal("Unexpected error routing request:", err)
	}
	if (*h).(MockHandler).Name != h1.Name {
		t.Error("Routed to wrong handler")
	}
}

type PanickingHandler struct{}

func (h PanickingHandler) HandleRequest(request Request) error {
	panic("I fail")
}

func TestPanicRecovery(t *testing.T) {
	server := NewServer()

	server.AddHandler("requesttype1.*", PanickingHandler{})

	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Server request handling panicked without recovering")
		}
	}()

	server.HandleRequest(MockRequest{Payload: "requesttype1 blah blah"})
}
