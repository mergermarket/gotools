package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

// InternalHealthCheck is used by our infrastructure to check the service is listening
func InternalHealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Healthy")
}

// InternalLogConfig creates an http handler which logs out the app's config
func InternalLogConfig(config interface{}, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Info(fmt.Sprintf("Application config - %+v", config))
		fmt.Fprint(w, "Logged the config")
	}
}

// InternalLRuntimeInfo logs current running version of GO
func InternalRuntimeInfo(w http.ResponseWriter, _ *http.Request) {
	versionInfo := struct {
		Platform string `json:"platform"`
		Version  string `json:"version"`
	}{
		"go",
		runtime.Version(),
	}
	versionInfoBytes, err := json.Marshal(versionInfo)
	if err != nil {
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(versionInfoBytes)
	if err != nil {
		fmt.Printf("Error writing versionInfo response: %v", err)
	}

}
