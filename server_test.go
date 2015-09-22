package dumbdumb

import (
	"testing"
)

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
