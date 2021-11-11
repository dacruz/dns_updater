package ipfy

import (
	"net"
	"net/http"
	"context"
	"testing"
	"time"
)

func TestFetchCurrentIp(t *testing.T) {
	server := startServer()
	defer stopServer(server)
	
	currectIp, _ := FetchCurrentIp("http://localhost:7000/current/ip")
	
	if ! net.ParseIP("10.0.0.1").Equal(currectIp) {
		t.Fatal("currectIp does not have the expected value")
	}
	
}

func TestFailToRequestCurrentIp(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentIp("WRONG://localhost:7000/current/ip")
	
	
	if err == nil {
		t.Fatal("FetchCurrentIp should have failed to execute GET")
	}
	
}

func TestFailToFetchCurrentIp(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentIp("http://localhost:7000/WRONG/PATH")
	
	if err == nil {
		t.Fatal("FetchCurrentIp should not succeed on non 2xx")
	}
	
}

func TestFailToParseFetchCurrentIpResponse(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentIp("http://localhost:7000/WRONG/IP")
	
	if err == nil {
		t.Fatal("FetchCurrentIp should not have returned a valid ip")
	}
	
}

func startServer() *http.Server {
	router := http.NewServeMux() 
	router.HandleFunc("/current/ip", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("10.0.0.1"))
	})

	router.HandleFunc("/WRONG/IP", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("WRONG_IP"))
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
