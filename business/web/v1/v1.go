//Package v1 contains all the interface use by this version
//of the middleware
package v1

//ErrorResponse is the form used by the API to send responses
//from failures
type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}
