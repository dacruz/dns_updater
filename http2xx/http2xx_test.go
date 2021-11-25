package http2xx

import (
	"io/ioutil"
	"net/http"
	"testing"
)

var handlers = map[string]func(http.ResponseWriter, *http.Request){
	"/echo/header": func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(r.Header.Get("Ping")))
	},
	"/echo/body": func(rw http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		rw.Write(body)
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

	headers := map[string]string{
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

func TestFailToExecutePut(t *testing.T) {
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	_, err := Put("WRONG://localhost:7000/current/ip", nil, nil)

	if err == nil {
		t.Fatal("Put should have failed to execute Put")
	}

}

func TestFailToPut2xx(t *testing.T) {
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	_, err := Put("http://localhost:7000/WRONG/PATH", nil, nil)

	if err == nil {
		t.Fatal("Put should not succeed on non 2xx")
	}

}

func TestPutWithHeaders(t *testing.T) {
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	headers := map[string]string{
		"Ping": "Pong",
	}

	resp, _ := Put("http://localhost:7000/echo/header", headers, nil)

	if string(resp) != "Pong" {
		t.Fatal("Put should include the headers")
	}

}

func TestPutWithoutHeaders(t *testing.T) {

	server := StartStubServer(handlers)
	defer StopStubServer(server)

	resp, _ := Put("http://localhost:7000/echo/header", nil, nil)

	if string(resp) != "" {
		t.Fatal("Put should not include the headers if it is empty")
	}

}

func TestPutWithBody(t *testing.T) {
	server := StartStubServer(handlers)
	defer StopStubServer(server)

	resp, _ := Put("http://localhost:7000/echo/body", nil, []byte("body"))

	if string(resp) != "body" {
		t.Fatal("Put should include the body")
	}

}

func TestPutWithoutBody(t *testing.T) {

	server := StartStubServer(handlers)
	defer StopStubServer(server)

	resp, _ := Put("http://localhost:7000/echo/body", nil, nil)

	if string(resp) != "" {
		t.Fatal("Put should not include the headers if it is empty")
	}

}
