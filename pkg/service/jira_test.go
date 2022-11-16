package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	c := Close{
		Issue: Issue{
			Key: "123",
			Fields: Fields{
				Status: Status{
					Name: "Done",
				},
				MrIid:           "MrIid",
				ProjectId:       "ProjectId",
				HotFix:          "HotFix",
				MergeRequestUrl: "MergeRequestUrl",
			},
		},
	}
	err := c.CloseValidate()
	assert.Nil(t, err)
}

func TestCloseProcessedData(t *testing.T) {
	c := Close{
		Issue: Issue{
			Key: "123",
			Fields: Fields{
				Status: Status{
					Name: "Done",
				},
				MrIid:           "MrIid",
				ProjectId:       "ProjectId",
				HotFix:          "HotFix",
				MergeRequestUrl: "MergeRequestUrl",
			},
		},
	}
	pcjk := ProduceJiraCloseKafka{}
	pcjk.ProcessedData(c)
	expectedPcjk := ProduceJiraCloseKafka{
		IssueKey:        "123",
		Status:          "Done",
		MrIid:           "MrIid",
		ProjectId:       "ProjectId",
		HotFix:          "HotFix",
		MergeRequestUrl: "MergeRequestUrl",
	}
	assert.Equal(t, expectedPcjk, pcjk, "They should be equal")
}
