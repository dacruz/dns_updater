package http2xx

import (
	"net/http"
	"testing"
)

var handlers = map[string]func(http.ResponseWriter, *http.Request) {
	"/echo/header": func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(r.Header.Get("Ping")))
	},
}

func TestFailToExecuteGet(t *testing.T) {
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	_, err := Get("WRONG://localhost:7000/current/ip", nil)
	
	
	if err == nil {
		t.Fatal("Get should have failed to execute GET")
	}
	
}

func TestFailToGet2xx(t *testing.T) {
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	_, err := Get("http://localhost:7000/WRONG/PATH", nil)
	
	if err == nil {
		t.Fatal("Get should not succeed on non 2xx")
	}
	
}

func TestGetWithHeaders(t *testing.T) {
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	headers := map[string]string {
		"Ping": "Pong",
	}
	
	resp, _ := Get("http://localhost:7000/echo/header", headers)

	if string(resp) != "Pong" {
		t.Fatal("Get should include the headers")
	}
	
}

func TestGetWithoutHeaders(t *testing.T) {
	
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	resp, _ := Get("http://localhost:7000/echo/header", nil)
	
	if string(resp) != "" {
		t.Fatal("Get should not include the headers if it is empty")
	}
	
}
