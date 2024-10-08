package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const environment = "PROD"

func getEnv(key string) string {
	if environment == "DEV" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		return os.Getenv(key)
	}

	return os.Getenv(key)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "applcation/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Enable CORS for each request
        enableCors(&w)

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content for preflight
            return
        }

        // Call the wrapped apiFunc
        if err := f(w, r); err != nil {
            WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
        }
    }
}
