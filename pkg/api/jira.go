package api

import (
	"fmt"
	"net/http"

	"github.com/0quz/gitlab-jira-cm/pkg/middleware"
	"github.com/0quz/gitlab-jira-cm/pkg/service"
)

// handling incoming requests from jira
func (ma MarcoAPI) HandleClose() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := ma.MarcoService.JiraClose(r.Context().Value(middleware.KeyClose{}).(service.Close))
		if err != nil {
			fmt.Println(err)
		}
	}
}
