package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Response(ctx context.Context, w http.ResponseWriter, statusCode int, data any) error {
	//@todo Set the response to the sentry log
	// Set the status code for the request logger middleware.
	if err := SetStatusCode(ctx, statusCode); err != nil {
		return err
	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}
	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")
	// Write the status code to the response.
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
