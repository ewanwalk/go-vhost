package vhost

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {

	router := New()
	// Test A
	// Checking routing based on "localhost" as the hostname
	router.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world a"))
	}, "localhost")

	// Test B
	// Checking routing based on "127.0.0.1" as the hostname
	router.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world b"))
	}, "127.0.0.1")

	assert := func(t *testing.T, s *httptest.Server, url, expectedBody string) {
		resp, err := s.Client().Get(url)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Fatalf("receieved non 200 status code: %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("unexpected error reading body %v", err)
		}

		if !bytes.Equal(body, []byte(expectedBody)) {
			t.Fatalf("expected response %s got %s", expectedBody, string(body))
		}
	}

	mock := func(url string) (*httptest.Server, error) {
		server := httptest.NewUnstartedServer(router)
		listen, err := net.Listen("tcp", url)
		if err != nil {
			return nil, err
		}

		// replace the listener
		server.Listener.Close()
		server.Listener = listen
		server.Start()

		return server, nil
	}

	t.Run("httpServerA", func(t *testing.T) {
		server, err := mock("127.0.0.1:12345")
		if err != nil {
			t.Fatalf("failed to mock server: %v", err)
		}

		defer server.Close()
		assert(t, server, "http://localhost:12345/eofmwefwfji", "hello world a")
	})

	t.Run("httpServerB", func(t *testing.T) {
		server, err := mock("127.0.0.1:12345")
		if err != nil {
			t.Fatalf("failed to mock server: %v", err)
		}

		defer server.Close()
		assert(t, server, "http://127.0.0.1:12345/eofmwefwfji", "hello world b")
	})
}

func TestRouter_stripStrict(t *testing.T) {

	router := New()

	url := "www.localhost"

	if res := router.stripStrict(url); res != "localhost" {
		t.Fatalf("expected %s got %s", "localhost", res)
	}

	router.Strict = true

	if res := router.stripStrict(url); res != "www.localhost" {
		t.Fatalf("expected %s got %s", "www.localhost", res)
	}

}
