package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0quz/gitlab-jira-cm/pkg/service"
)

// the key for context after the middleware section.
type KeyMergeRequest struct{}

// the middleware that checks to convert data from JSON and validates
func MergeRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			mr := service.MergeRequest{}
			err := mr.FromJSON(r.Body)
			if err != nil {
				http.Error(rw, fmt.Sprintf("Error reading gitlab merge request: %s", err), http.StatusBadRequest)
				return
			}
			err = mr.MergeRequestValidate()
			if err != nil {
				http.Error(rw, fmt.Sprintf("Error validating gitlab merge request: %s", err), http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), KeyMergeRequest{}, mr)
			r = r.WithContext(ctx)

		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
