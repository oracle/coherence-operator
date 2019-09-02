package fakes

import (
	"encoding/json"
	"fmt"
	"github.com/oracle/coherence-operator/test/e2e/helper"
	"net/http"
)

// FakeRestServer can be use to server ReST requests in tests.
type FakeRestServer interface {
	// Add a handler to respond to a query for a specific path.
	AddHandler(path string, handler http.Handler)
	// Add a response value to send in response to a query for a specific path.
	// Primitive values will be sent to the response stream as-is whereas structs, maps, slices
	// etc, will be marshalled into json and sent to the response.
	AddResponse(path string, value interface{})
	// Close the web server
	Close()
	// Get the port that the server is bound to
	GetPort() int32
	// Get a url to use to perform a query for the specified path.
	// The returned value will be in the form http://127.0.0.1:<port>/<path>
	// where portis the server listen port and path is the specified path
	GetURL(path string) string
}

// NewFakeRestServer creates and starts a fake ReST server.
func NewFakeRestServer() (FakeRestServer, error) {
	f := &fakeServer{}
	ports := helper.GetAvailablePorts()
	port, err := ports.Next()
	if err != nil {
		return nil, err
	}

	svr := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Handler: f}
	f.server = svr
	f.port = port
	f.responses = make(map[string]http.Handler)

	go func() {
		err := svr.ListenAndServe()
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	return f, nil
}

// ----- internal -----------------------------------------------------------

type fakeServer struct {
	responses map[string]http.Handler
	server    *http.Server
	port      int32
}

func (f *fakeServer) GetPort() int32 {
	return f.port
}

func (f *fakeServer) GetURL(path string) string {
	sep := "/"
	if path[0] == '/' {
		sep = ""
	}
	return fmt.Sprintf("http://127.0.0.1:%d%s%s", f.port, sep, path)
}

func (f *fakeServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := f.responses[req.URL.Path]
	if handler == nil {
		w.WriteHeader(404)
		_, _ = fmt.Fprint(w, "not found")
	} else {
		handler.ServeHTTP(w, req)
	}
}

func (f *fakeServer) AddResponse(path string, resp interface{}) {
	handler := &response{value: resp}
	f.AddHandler(path, handler)
}

func (f *fakeServer) AddHandler(path string, handler http.Handler) {
	f.responses[path] = handler
}

func (f *fakeServer) Close() {
	_ = f.server.Close()
}

type response struct {
	value interface{}
}

var _ http.Handler = &response{}

func (r *response) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch v := r.value.(type) {
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128:
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, "%v", v)
	case string:
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, "%s", r.value)
	default:
		data, err := json.Marshal(r.value)
		if err != nil {
			w.WriteHeader(500)
			_, _ = fmt.Fprintf(w, "Error marshalling response to json.\n%s", err.Error())
		} else {
			w.WriteHeader(200)
			_, err = w.Write(data)
			if err != nil {
				fmt.Print(err)
			}
		}
	}
}
