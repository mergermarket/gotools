package httputil

import "net/http"
import "errors"
import (
	"fmt"
	"strings"
)

// Validate verifies that all strings passed as `params` are present as keys in req.Header.
// An error listing the missing params is returned if at least one param is missing.
func Validate(params []string, req *http.Request) error {
	missing := getMissingKeys(req, params)
	if len(missing) == 0 {
		return nil
	}
	return errors.New("The following required params are missing from the request (as headers) : " + strings.Join(missing, ","))
}

// ValidateParamsHandler wraps a http.Handler with validation to check if the specified params are
// supplied as HTTP headers.
// If not, a 400 is returned with a list of the missing params.
// If all required params supplied, the supplied handler is called transparently.
func ValidateParamsHandler(h http.Handler, params ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := Validate(params, req)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%s", err.Error())
			return
		}
		h.ServeHTTP(w, req)
	})
}

func getMissingKeys(req *http.Request, required []string) []string {
	missing := []string{}
	for _, param := range required {
		if req.Header.Get(param) == "" {
			missing = append(missing, param)
		}
	}
	return missing
}
