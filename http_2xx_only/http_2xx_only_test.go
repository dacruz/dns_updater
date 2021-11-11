package http_2xx_only

import (
	"net/http"
	"context"
	"testing"
	"time"
)

func TestFailToExecuteGet(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := Get("WRONG://localhost:7000/current/ip", nil)
	
	
	if err == nil {
		t.Fatal("Get should have failed to execute GET")
	}
	
}

func TestFailToGet2xx(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := Get("http://localhost:7000/WRONG/PATH", nil)
	
	if err == nil {
		t.Fatal("Get should not succeed on non 2xx")
	}
	
}


func TestGetWithHeaders(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	headers := map[string]string {
		"Ping": "Pong",
	}
	
	resp, _ := Get("http://localhost:7000/echo/header", headers)

	if string(resp) != "Pong" {
		t.Fatal("Get should include the headers")
	}
	
}

func TestGetWithoutHeaders(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	resp, _ := Get("http://localhost:7000//echo/header", nil)
	
	if string(resp) != "" {
		t.Fatal("Get should not include the headers if it is empty")
	}
	
}

func startServer() *http.Server {
	router := http.NewServeMux() 

	router.HandleFunc("/echo/header", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(r.Header.Get("Ping")))
	})

	server := &http.Server{
		Addr:         "localhost:7000",
		Handler:      router,
	}

	go server.ListenAndServe()
	
	return server
}

func stopServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3* time.Second)
  	defer cancel()

  	server.SetKeepAlivesEnabled(false)
  	if err := server.Shutdown(ctx); err != nil {
    	return err
  	}
	
	return server.Shutdown(ctx)
}