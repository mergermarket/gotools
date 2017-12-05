package tools

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInternalHealthcheckRouter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(InternalHealthCheck))
	response, err := http.Get(server.URL + "/internal/healthcheck")

	if err != nil {
		t.Fatal("I got an error requesting healthcheck ", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("Expected a 200 but I got ", response.StatusCode)
	}
}

func TestInternalLRuntimeInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(InternalRuntimeInfo))
	response, err := http.Get(server.URL + "/internal/runtime-info")

	if err != nil {
		t.Fatal("I got an error requesting runtime-info ", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("Expected a 200 but I got ", response.StatusCode)
	}
}
