package godaddy

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestFetchCurrentRecordValue(t *testing.T) {
	server := startServer()
	defer stopServer(server)
	
	currectIp, _ := FetchCurrentRecordValue("http://localhost:7000/v1", "poiuytre.nl", "@", "API_KEY")
	
	if ! net.ParseIP("10.0.0.1").Equal(currectIp) {
		t.Fatal("record does not have the expected value")
	}
	
}

func TestFailToRequestCurrentRecordValue(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentRecordValue("WRONG://localhost:7000/v1", "poiuytre.nl", "@", "API_KEY")
	
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should have failed to execute GET")
	}
	
}

func TestFailToFetchCurrentRecordValue(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/WRONG/v2", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should not succeed on non 2xx")
	}
	
}

func TestFailToParseFetchCurrentRecordValueResponse(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/v1/WRONG/IP", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should not have returned a valid ip")
	}
	
}

func startServer() *http.Server {
	router := http.NewServeMux() 

	router.HandleFunc("/v1/domains/poiuytre.nl/records/A/@", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "sso-key API_KEY" {
			rw.Write([]byte("10.0.0.1"))
		} else {
			rw.Write([]byte("sso-key motherf***er, do you speak it?!"))
		}
	})

	router.HandleFunc("/v1/WRONG/IP/domains/poiuytre.nl/records/A/@", func(rw http.ResponseWriter, r *http.Request) {
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