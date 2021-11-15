package http2xx

import (
	"context"
	"net/http"
	"time"
)

func StartStubServer(handlers map[string]func(http.ResponseWriter, *http.Request)) *http.Server {
	router := http.NewServeMux()

	for pattern, handler := range handlers {
		router.HandleFunc(pattern, handler)
	}

	server := &http.Server{
		Addr:    "localhost:7000",
		Handler: router,
	}

	go server.ListenAndServe()
	time.Sleep(1 * time.Millisecond)

	return server
}

func StopStubServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	server.Shutdown(ctx)

	return server.Shutdown(ctx)
}
