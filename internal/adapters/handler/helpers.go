package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("decode json: %v", err)
		return v, friendlyJSONError(err)
	}
	return v, nil
}

func friendlyJSONError(err error) error {
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return fmt.Errorf("field '%s' has invalid type: expected %s", typeErr.Field, typeErr.Type)
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return fmt.Errorf("malformed JSON at position %d", syntaxErr.Offset)
	}

	return fmt.Errorf("invalid request body")
}
