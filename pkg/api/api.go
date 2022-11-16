package api

import (
	"net/http"

	"github.com/0quz/gitlab-jira-cm/pkg/service"
)

// dependency part
type MarcoAPI struct {
	MarcoService service.MarcoService
}

func NewMarcoAPI(m service.MarcoService) MarcoAPI {
	return MarcoAPI{MarcoService: m}
}

// basic health check
func (ma MarcoAPI) HandleHealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			rw.Write([]byte("ok"))
			return
		}
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
