package api

import "net/http"

type settingsResponse struct {
	ServerVersion string `json:"server_version"`
	TokenMasked   string `json:"token_masked"`
}

func settingsHandler(token string) http.HandlerFunc {
	masked := "not set"
	if token != "" {
		if len(token) <= 4 {
			masked = "****"
		} else {
			masked = token[:2] + "••••" + token[len(token)-2:]
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, settingsResponse{
			ServerVersion: "0.1.0",
			TokenMasked:   masked,
		})
	}
}
