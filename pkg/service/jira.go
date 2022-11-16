package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Struct for incoming request
type Close struct {
	Issue Issue `json:"issue" validate:"required"`
}

type Issue struct {
	Key    string `json:"key" validate:"required"`
	Fields Fields `json:"fields" validate:"required"`
}

type Fields struct {
	Status          Status         `json:"status" validate:"required"`
	AffectedDomain  AffectedDomain `json:"customfield_38460"`
	ChangePriorty   ChangePriorty  `json:"customfield_36061"`
	MrIid           string         `json:"customfield_38368" validate:"required"`
	ProjectId       string         `json:"customfield_38372" validate:"required"`
	Summary         string         `json:"summary"`
	HotFix          string         `json:"customfield_38371" validate:"required"`
	MergeRequestUrl string         `json:"customfield_38373" validate:"required"`
}

type Status struct {
	Name string `json:"name" validate:"required,name"`
}

type AffectedDomain struct {
	Department string `json:"value"`
	Child      Child  `json:"child"`
}

type Child struct {
	Team string `json:"value"`
}

type ChangePriorty struct {
	ChangePriortyV string `json:"value"`
}

// Struct for Kafka feed
type ProduceJiraCloseKafka struct {
	IssueKey        string `json:"issue_key"`
	Status          string `json:"status"`
	MrIid           string `json:"mr_iid"`
	ProjectId       string `json:"project_id"`
	Department      string `json:"department"`
	Team            string `json:"team"`
	Summary         string `json:"summary"`
	HotFix          string `json:"hot_fix"`
	MergeRequestUrl string `json:"merge_url"`
	ChangePriority  string `json:"change_priority"`
}

// JSON to Struct converter for incoming request
func (c *Close) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(c)
}

// Prepare requested data
func (pjck *ProduceJiraCloseKafka) ProcessedData(c Close) {
	pjck.IssueKey = c.Issue.Key
	pjck.Status = c.Issue.Fields.Status.Name
	pjck.MrIid = c.Issue.Fields.MrIid
	pjck.ProjectId = c.Issue.Fields.ProjectId
	pjck.Department = c.Issue.Fields.AffectedDomain.Department
	pjck.Team = c.Issue.Fields.AffectedDomain.Child.Team
	pjck.Summary = c.Issue.Fields.Summary
	pjck.HotFix = c.Issue.Fields.HotFix
	pjck.MergeRequestUrl = c.Issue.Fields.MergeRequestUrl
	pjck.ChangePriority = c.Issue.Fields.ChangePriorty.ChangePriortyV
}

// Jira close handler
// It feeds the kafka topic.
func (ms *MarcoService) JiraClose(c Close) error {
	pjck := ProduceJiraCloseKafka{}
	pjck.ProcessedData(c)
	p, err := json.Marshal(pjck)
	if err != nil {
		return fmt.Errorf("JiraClose: ProcessedData to marshal error: %w", err)
	}
	err = ProduceKafkaTopic(os.Getenv("topic_jira_close"), p)
	if err != nil {
		return fmt.Errorf("ProduceKafkaTopic: Produce kafka topic: %s error: %w", os.Getenv("topic_jira_close"), err)
	}
	//fmt.Printf(string(p))
	return nil
}
