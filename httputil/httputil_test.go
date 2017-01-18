package httputil

import "testing"
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func TestValidate_ReturnsErrorWhenNoRequiredParamsSupplied(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("Content-Type", "application/json")

	err := Validate([]string{"Required-One", "Required-Two"}, req)

	notAsExpected := err == nil ||
		!strings.Contains(err.Error(), "Required-One") ||
		!strings.Contains(err.Error(), "Required-Two")
	if notAsExpected {
		t.Error("Expected error listing Required-One and Required-Two")
	}
}

func TestValidate_ReturnsNoErrorWhenAllRequiredParamsSupplied(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("Required-Two", "something")
	req.Header.Add("Required-One", "something-else")
	req.Header.Add("Other-Header", "not-relevant")

	err := Validate([]string{"Required-One", "Required-Two"}, req)

	if err != nil {
		t.Error("Expected no errors as both required params are present")
	}

}

type TestHandler struct{}

func (h TestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello from TestHandler")
}

func TestValidateParamsHandler_RequiredParamsPassed(t *testing.T) {
	testHandler := TestHandler{}
	testHandlerWithValidation := ValidateParamsHandler(testHandler, "Required-One", "Required-Two")

	ts := httptest.NewServer(testHandlerWithValidation)
	defer ts.Close()

	client := &http.Client{}
	req := httptest.NewRequest("GET", ts.URL, nil)
	req.RequestURI = ""
	req.Header.Add("Required-Two", "something")
	req.Header.Add("Required-One", "something-else")
	req.Header.Add("Other-Header", "not-relevant")

	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200 but got %d", resp.StatusCode)
	}
}

func TestValidateParamsHandler_SomeRequiredParamsMissing(t *testing.T) {
	testHandler := TestHandler{}
	testHandlerWithValidation := ValidateParamsHandler(testHandler, "Required-One", "Required-Two", "Required-Three")

	ts := httptest.NewServer(testHandlerWithValidation)
	defer ts.Close()

	client := &http.Client{}
	req := httptest.NewRequest("GET", ts.URL, nil)
	req.RequestURI = ""
	req.Header.Add("Required-Two", "something")
	req.Header.Add("Required-One", "something-else")
	req.Header.Add("Other-Header", "not-relevant")

	resp, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("Expected status code 400 but got %d", resp.StatusCode)
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	body := string(bodybytes)
	fmt.Println(body)
	if !strings.Contains(body, "Required-Three") {
		t.Error("Should have error message with 'Required-Three' in the response body")
	}
	if strings.Contains(body, "Required-One") || strings.Contains(body, "Required-Two") {
		t.Error("Should not have 'Required-One' or 'Required-Two' in the response body")
	}
}
