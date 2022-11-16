package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/0quz/gitlab-jira-integration/pkg/middleware"
	"github.com/0quz/gitlab-jira-integration/pkg/service"
	"gorm.io/gorm"
)

// handling incoming requests from gitlab
func (ma MarcoAPI) HandleMergeRequest() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mr := r.Context().Value(middleware.KeyMergeRequest{}).(service.MergeRequest)
		err := ma.MarcoService.FindEventByMrUrl(mr.ObjectAttributes.Url)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if strings.Contains(mr.ObjectAttributes.SourceBranch, "hot-fix") {
					go ma.MarcoService.HotFix(mr)
				} else {
					go ma.MarcoService.Standard(mr)
				}
			} else {
				fmt.Println("HandleMergeRequest: FindEventByMrUrl db error: ", err.Error())
			}
		}
	}
}
