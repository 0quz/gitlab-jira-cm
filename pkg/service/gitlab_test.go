package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeRequestData(t *testing.T) {
	mr := MergeRequest{
		Project: Project{
			Id:        2808,
			Name:      "project",
			Namespace: "group",
		},
		ObjectAttributes: ObjectAttributes{
			Url:          "https://gitlab.com/group/project/-/merge_requests/23",
			Iid:          23,
			SourceBranch: "CM-599512-3081117",
			TargetBranch: "prod-release",
			State:        "opened",
			Action:       "open",
			Title:        "ok",
		},
	}
	err := mr.MergeRequestValidate()
	assert.Nil(t, err)
}

func TestMergeRequestProcessedData(t *testing.T) {
	mr := MergeRequest{
		Project: Project{
			Id:        2808,
			Name:      "project",
			Namespace: "group",
		},
		ObjectAttributes: ObjectAttributes{
			LastCommit: LastCommit{
				Author: Author{
					Email: "username@gmail.com",
				},
			},
			Url:          "https://gitlab.com/group/project/-/merge_requests/23",
			Iid:          23,
			TargetBranch: "prod-release",
			State:        "opened",
			Action:       "open",
		},
	}
	pgk := ProduceGitlabKafka{}
	pgk.ProcessedData(mr)
	expectedPgk := ProduceGitlabKafka{
		User: KUser{
			Name:     "John Doe",
			Username: "jdoe",
		},
		Project: KProject{
			ProjectId:   "2808",
			ProjectName: "project",
			Namespace:   "group",
		},
		MergeRequest: KMergeRequest{
			Iid: "23",
			Url: "https://gitlab.com/group/project/-/merge_requests/23",
		},
		Watchers: pgk.Watchers,
		Assignee: pgk.Assignee,
	}
	assert.Equal(t, expectedPgk, pgk, "They should be equal")
}
