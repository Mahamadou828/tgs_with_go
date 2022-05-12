package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Response(ctx context.Context, w http.ResponseWriter, statusCode int, data any) *RequestError {
	//@todo Set the response to the sentry log
	// Set the status code for the request logger middleware.
	if err := SetStatusCode(ctx, statusCode); err != nil {
		return NewRequestError(
			NewShutdownError("can't set status code inside the header"),
			http.StatusInternalServerError,
		)
	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		return NewRequestError(
			NewShutdownError("can't parse json response: "+err.Error()),
			http.StatusInternalServerError,
		)
	}
	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")
	// Write the status code to the response.
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return NewRequestError(
			NewShutdownError("can't write json response: "+err.Error()),
			http.StatusInternalServerError,
		)
	}

	return nil
}
