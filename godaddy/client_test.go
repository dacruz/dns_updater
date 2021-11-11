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

func TestFailToParseFetchCurrentRecordValueResponse(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/v1/WRONG/RESPONSE", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should not have returned a valid json")
	}
	
}

func TestFailToParseFetchCurrentRecordValue(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/v1/WRONG/IP", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should not have returned a valid ip")
	}
	
}

func TestFailFetchCurrentRecordValueOnNon2xx(t *testing.T) {
	server := startServer()
	defer stopServer(server)

	_, err := FetchCurrentRecordValue("http://localhost:7000/NOT_2XX", "poiuytre.nl", "@", "API_KEY")
	
	if err == nil {
		t.Fatal("FetchCurrentRecordValue should fail on non 2xx")
	}
	
}

func startServer() *http.Server {
	router := http.NewServeMux() 

	router.HandleFunc("/v1/domains/poiuytre.nl/records/A/@", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "sso-key API_KEY" {
			rw.Write([]byte(`[{"data":"10.0.0.1","name":"@","ttl":600,"type":"A"}]`))
		} else {
			rw.Write([]byte("sso-key motherf***er, do you speak it?!"))
		}
	})

	router.HandleFunc("/v1/WRONG/RESPONSE/domains/poiuytre.nl/records/A/@", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("WRONG_IP"))
	})

	router.HandleFunc("/v1/WRONG/IP/domains/poiuytre.nl/records/A/@", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(`[{"data":"10","name":"@","ttl":600,"type":"A"}]`))
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
