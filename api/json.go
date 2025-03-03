package api

import (
    "net/http"
    "encoding/json"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

// write a json error with the given message on the body
func writeJSONError(w http.ResponseWriter, status int, message string) error {
    type envelope struct{
	Error string `json:"error"`
    }

    return writeJSON(w, status, &envelope{Error: message})
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
    max_bytes := 1_048_578

    r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))

    dec := json.NewDecoder(r.Body)
    //strict mode
    dec.DisallowUnknownFields()

    return dec.Decode(data)
}
