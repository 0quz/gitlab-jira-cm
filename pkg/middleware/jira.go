package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/0quz/gitlab-jira-cm/pkg/service"
)

// the key for context after the middleware section.
type KeyClose struct{}

// the middleware that checks to convert data from JSON and validates
func CloseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cd := service.Close{}
			err := cd.FromJSON(r.Body)
			if err != nil {
				http.Error(rw, fmt.Sprintf("Error reading jira close: %s", err), http.StatusBadRequest)
				return
			}
			err = cd.CloseValidate()
			if err != nil {
				http.Error(rw, fmt.Sprintf("Error validating jira close: %s", err), http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), KeyClose{}, cd)
			r = r.WithContext(ctx)

		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
